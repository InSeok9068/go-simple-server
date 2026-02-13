// go build -o ./scaffold.exe ./cmd/scaffold
// go run ./cmd/scaffold

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

const (
	deployConfigPath = ".linux/scaffold/deploy-config.json"
	caddyFilePath    = ".linux/caddy/Caddyfile"
	serviceDirPath   = ".linux/systemctl"
	customProjectOpt = "[직접 입력]"
)

var (
	projectNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
)

type DeployConfig struct {
	SSHHost            string `json:"ssh_host"`
	SSHUser            string `json:"ssh_user"`
	SSHPort            int    `json:"ssh_port"`
	SSHKeyPath         string `json:"ssh_key_path"`
	RemoteAppDir       string `json:"remote_app_dir"`
	RemoteSystemdDir   string `json:"remote_systemd_dir"`
	RemoteCaddyfile    string `json:"remote_caddyfile"`
	RemoteSrvRoot      string `json:"remote_srv_root"`
	SkipSSHKey         bool   `json:"skip_ssh_key"`
	DisableHostKeyTest bool   `json:"disable_host_key_test"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var scaffoldType string
	if err := survey.AskOne(&survey.Select{
		Message: "작업할 내용을 선택하세요.",
		Options: []string{
			"프로젝트 배포",
			"프로젝트 배포 회수",
		},
	}, &scaffoldType); err != nil {
		return err
	}

	switch scaffoldType {
	case "프로젝트 배포":
		return runDeployProject()
	case "프로젝트 배포 회수":
		return runUndeployProject()
	default:
		return fmt.Errorf("알 수 없는 작업 유형입니다: %s", scaffoldType)
	}
}

func runDeployProject() error {
	logTitle("프로젝트 첫 배포 시작")
	projectName, err := askProjectName("첫 배포할 프로젝트를 선택하세요.")
	if err != nil {
		return err
	}

	localBinaryPath, err := resolveLocalBinaryPath(projectName)
	if err != nil {
		return err
	}

	cfg, err := askDeployConfig()
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("배포 요약")
	fmt.Printf("- 프로젝트: %s\n", projectName)
	fmt.Printf("- 서버: %s@%s:%d\n", cfg.SSHUser, cfg.SSHHost, cfg.SSHPort)
	fmt.Printf("- 로컬 바이너리: %s\n", localBinaryPath)
	fmt.Printf("- 원격 바이너리 경로: %s\n", path.Join(cfg.RemoteAppDir, projectName))
	fmt.Printf("- 원격 서비스 파일: %s\n", path.Join(cfg.RemoteSystemdDir, projectName+".service"))
	fmt.Printf("- 원격 Caddyfile: %s (로컬 파일 덮어쓰기)\n", cfg.RemoteCaddyfile)
	fmt.Println("")

	var ok bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "첫 배포 자동화를 실행하시겠습니까? ",
		Default: false,
	}, &ok, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if !ok {
		fmt.Println("배포를 취소했습니다.")
		return nil
	}

	logStep("필수 명령어 확인")
	if err := ensureCommands("ssh", "scp"); err != nil {
		return err
	}
	return deployProject(projectName, localBinaryPath, cfg)
}

func runUndeployProject() error {
	logTitle("프로젝트 회수 시작")
	projectName, err := askProjectName("회수할 프로젝트를 선택하세요.")
	if err != nil {
		return err
	}

	cfg, err := askDeployConfig()
	if err != nil {
		return err
	}

	var removeSrvData bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("%s 서버 데이터 디렉터리(/srv 기준)도 삭제할까요? ", projectName),
		Default: true,
	}, &removeSrvData, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var removeLocalService bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "로컬 .linux/systemctl 서비스 파일도 제거할까요? ",
		Default: true,
	}, &removeLocalService, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("회수 요약")
	fmt.Printf("- 프로젝트: %s\n", projectName)
	fmt.Printf("- 서버: %s@%s:%d\n", cfg.SSHUser, cfg.SSHHost, cfg.SSHPort)
	fmt.Printf("- 원격 바이너리 경로: %s\n", path.Join(cfg.RemoteAppDir, projectName))
	fmt.Printf("- 원격 서비스 파일: %s\n", path.Join(cfg.RemoteSystemdDir, projectName+".service"))
	fmt.Printf("- 원격 Caddyfile: %s (로컬 파일 덮어쓰기)\n", cfg.RemoteCaddyfile)
	fmt.Printf("- /srv/%s 삭제 여부: %t\n", projectName, removeSrvData)
	fmt.Println("")

	var ok bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "프로젝트 회수 자동화를 실행하시겠습니까? ",
		Default: false,
	}, &ok, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if !ok {
		fmt.Println("회수를 취소했습니다.")
		return nil
	}

	logStep("필수 명령어 확인")
	if err := ensureCommands("ssh", "scp"); err != nil {
		return err
	}
	return undeployProject(projectName, cfg, removeSrvData, removeLocalService)
}

func deployProject(projectName, localBinaryPath string, cfg DeployConfig) error {
	logStep("로컬 서비스 디렉터리 준비")
	if err := os.MkdirAll(serviceDirPath, 0755); err != nil {
		return err
	}

	servicePath := filepath.Join(serviceDirPath, projectName+".service")
	logStep("systemd 서비스 파일 생성")
	if err := writeServiceFile(servicePath, projectName, cfg.RemoteAppDir); err != nil {
		return err
	}

	remoteBinaryTmp := fmt.Sprintf("/tmp/%s", projectName)
	remoteServiceTmp := fmt.Sprintf("/tmp/%s.service", projectName)
	remoteCaddyTmp := fmt.Sprintf("/tmp/Caddyfile.%s", projectName)
	logStep("바이너리 업로드")
	if err := scpFile(cfg, localBinaryPath, remoteBinaryTmp); err != nil {
		return err
	}
	logStep("서비스 파일 업로드")
	if err := scpFile(cfg, servicePath, remoteServiceTmp); err != nil {
		return err
	}
	logStep("Caddyfile 업로드")
	if err := scpFile(cfg, caddyFilePath, remoteCaddyTmp); err != nil {
		return err
	}

	script := renderDeployScript(projectName, cfg, remoteBinaryTmp, remoteServiceTmp, remoteCaddyTmp)
	logStep("원격 배포 스크립트 실행")
	if err := runRemoteScript(cfg, projectName+"-deploy", script); err != nil {
		return err
	}

	logSuccess("첫 배포 자동화가 완료되었습니다.")
	fmt.Println("첫 배포 자동화가 완료되었습니다.")
	return nil
}

func undeployProject(projectName string, cfg DeployConfig, removeSrvData, removeLocalService bool) error {
	servicePath := filepath.Join(serviceDirPath, projectName+".service")
	if removeLocalService {
		logStep("로컬 서비스 파일 정리")
		if err := os.Remove(servicePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	remoteCaddyTmp := fmt.Sprintf("/tmp/Caddyfile.%s", projectName)
	logStep("Caddyfile 업로드")
	if err := scpFile(cfg, caddyFilePath, remoteCaddyTmp); err != nil {
		return err
	}

	script := renderUndeployScript(projectName, cfg, remoteCaddyTmp, removeSrvData)
	logStep("원격 회수 스크립트 실행")
	if err := runRemoteScript(cfg, projectName+"-undeploy", script); err != nil {
		return err
	}

	logSuccess("프로젝트 회수 자동화가 완료되었습니다.")
	fmt.Println("프로젝트 회수 자동화가 완료되었습니다.")
	return nil
}

func askProjectName(message string) (string, error) {
	projectNames, err := listProjectNames()
	if err != nil {
		return "", err
	}
	projectNames = append(projectNames, customProjectOpt)

	var selected string
	if err := survey.AskOne(&survey.Select{
		Message: message,
		Options: projectNames,
	}, &selected, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}

	if selected == customProjectOpt {
		if err := survey.AskOne(&survey.Input{
			Message: "프로젝트 명을 직접 입력하세요: ",
		}, &selected,
			survey.WithValidator(survey.Required),
			survey.WithValidator(projectNameValidator)); err != nil {
			return "", err
		}
	}

	selected = strings.TrimSpace(selected)
	if !projectNamePattern.MatchString(selected) {
		return "", fmt.Errorf("프로젝트 명은 소문자/숫자/하이픈만 사용할 수 있습니다: %s", selected)
	}
	return selected, nil
}

func askDeployConfig() (DeployConfig, error) {
	cfg, exists, err := loadDeployConfig()
	if err != nil {
		return DeployConfig{}, err
	}

	useSaved := true
	if exists {
		if err := survey.AskOne(&survey.Confirm{
			Message: "저장된 서버 배포 설정을 사용할까요? ",
			Default: true,
		}, &useSaved, survey.WithValidator(survey.Required)); err != nil {
			return DeployConfig{}, err
		}
	}

	if useSaved && exists {
		return cfg, nil
	}

	var sshHost string
	var sshUser string
	var sshPortRaw string
	var sshKeyPath string
	var remoteAppDir string
	var remoteSystemdDir string
	var remoteCaddyfile string
	var remoteSrvRoot string
	var skipSSHKey bool
	var disableHostKeyTest bool

	if err := survey.AskOne(&survey.Input{
		Message: "서버 호스트(IP 또는 도메인): ",
		Default: cfg.SSHHost,
	}, &sshHost, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if err := survey.AskOne(&survey.Input{
		Message: "서버 사용자명: ",
		Default: cfg.SSHUser,
	}, &sshUser, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if err := survey.AskOne(&survey.Input{
		Message: "SSH 포트: ",
		Default: strconv.Itoa(cfg.SSHPort),
	}, &sshPortRaw,
		survey.WithValidator(survey.Required),
		survey.WithValidator(portValidator)); err != nil {
		return DeployConfig{}, err
	}
	sshPort, _ := strconv.Atoi(strings.TrimSpace(sshPortRaw))

	if err := survey.AskOne(&survey.Confirm{
		Message: "SSH 키 옵션(-i)을 생략하고 ssh-agent 인증을 사용할까요? ",
		Default: cfg.SkipSSHKey,
	}, &skipSSHKey, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if !skipSSHKey {
		if err := survey.AskOne(&survey.Input{
			Message: "SSH 개인키 경로: ",
			Default: cfg.SSHKeyPath,
		}, &sshKeyPath, survey.WithValidator(survey.Required)); err != nil {
			return DeployConfig{}, err
		}

		sshKeyPath = strings.TrimSpace(sshKeyPath)
		sshKeyPath, err = normalizeKeyPath(sshKeyPath)
		if err != nil {
			return DeployConfig{}, err
		}
	}

	if err := survey.AskOne(&survey.Input{
		Message: "원격 앱 디렉터리: ",
		Default: cfg.RemoteAppDir,
	}, &remoteAppDir, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if err := survey.AskOne(&survey.Input{
		Message: "원격 systemd 디렉터리: ",
		Default: cfg.RemoteSystemdDir,
	}, &remoteSystemdDir, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if err := survey.AskOne(&survey.Input{
		Message: "원격 Caddyfile 경로: ",
		Default: cfg.RemoteCaddyfile,
	}, &remoteCaddyfile, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if err := survey.AskOne(&survey.Input{
		Message: "원격 /srv 루트 경로: ",
		Default: cfg.RemoteSrvRoot,
	}, &remoteSrvRoot, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if err := survey.AskOne(&survey.Confirm{
		Message: "SSH StrictHostKeyChecking 검증을 끌까요? ",
		Default: true,
	}, &disableHostKeyTest, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	result := DeployConfig{
		SSHHost:            strings.TrimSpace(sshHost),
		SSHUser:            strings.TrimSpace(sshUser),
		SSHPort:            sshPort,
		SSHKeyPath:         strings.TrimSpace(sshKeyPath),
		RemoteAppDir:       strings.TrimSpace(remoteAppDir),
		RemoteSystemdDir:   strings.TrimSpace(remoteSystemdDir),
		RemoteCaddyfile:    strings.TrimSpace(remoteCaddyfile),
		RemoteSrvRoot:      strings.TrimSpace(remoteSrvRoot),
		SkipSSHKey:         skipSSHKey,
		DisableHostKeyTest: disableHostKeyTest,
	}

	var save bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "이 배포 설정을 로컬 파일에 저장할까요? ",
		Default: true,
	}, &save, survey.WithValidator(survey.Required)); err != nil {
		return DeployConfig{}, err
	}

	if save {
		if err := saveDeployConfig(result); err != nil {
			return DeployConfig{}, err
		}
		fmt.Printf("배포 설정을 저장했습니다: %s\n", deployConfigPath)
	}

	return result, nil
}

func loadDeployConfig() (DeployConfig, bool, error) {
	defaultCfg := defaultDeployConfig()

	if _, err := os.Stat(deployConfigPath); errors.Is(err, os.ErrNotExist) {
		return defaultCfg, false, nil
	}

	data, err := os.ReadFile(deployConfigPath)
	if err != nil {
		return DeployConfig{}, false, err
	}

	var cfg DeployConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DeployConfig{}, false, err
	}

	cfg = mergeDeployConfig(defaultCfg, cfg)
	if !cfg.SkipSSHKey {
		cfg.SSHKeyPath, err = normalizeKeyPath(cfg.SSHKeyPath)
		if err != nil {
			return DeployConfig{}, false, err
		}
	}

	return cfg, true, nil
}

func saveDeployConfig(cfg DeployConfig) error {
	if err := os.MkdirAll(filepath.Dir(deployConfigPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	return os.WriteFile(deployConfigPath, data, 0600)
}

func defaultDeployConfig() DeployConfig {
	homeDir, _ := os.UserHomeDir()
	return DeployConfig{
		SSHHost:            "",
		SSHUser:            "ubuntu",
		SSHPort:            22,
		SSHKeyPath:         filepath.Join(homeDir, ".ssh", "id_rsa"),
		RemoteAppDir:       "/home/ubuntu/app",
		RemoteSystemdDir:   "/etc/systemd/system",
		RemoteCaddyfile:    "/etc/caddy/Caddyfile",
		RemoteSrvRoot:      "/srv",
		SkipSSHKey:         false,
		DisableHostKeyTest: true,
	}
}

func mergeDeployConfig(base, override DeployConfig) DeployConfig {
	result := base
	if override.SSHHost != "" {
		result.SSHHost = override.SSHHost
	}
	if override.SSHUser != "" {
		result.SSHUser = override.SSHUser
	}
	if override.SSHPort > 0 {
		result.SSHPort = override.SSHPort
	}
	if override.SSHKeyPath != "" {
		result.SSHKeyPath = override.SSHKeyPath
	}
	if override.RemoteAppDir != "" {
		result.RemoteAppDir = override.RemoteAppDir
	}
	if override.RemoteSystemdDir != "" {
		result.RemoteSystemdDir = override.RemoteSystemdDir
	}
	if override.RemoteCaddyfile != "" {
		result.RemoteCaddyfile = override.RemoteCaddyfile
	}
	if override.RemoteSrvRoot != "" {
		result.RemoteSrvRoot = override.RemoteSrvRoot
	}
	result.SkipSSHKey = override.SkipSSHKey
	result.DisableHostKeyTest = override.DisableHostKeyTest
	return result
}

func normalizeKeyPath(keyPath string) (string, error) {
	if strings.TrimSpace(keyPath) == "" {
		return "", nil
	}
	if strings.HasPrefix(keyPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		keyPath = filepath.Join(homeDir, strings.TrimPrefix(keyPath, "~"))
	}
	return filepath.Clean(keyPath), nil
}

func listProjectNames() ([]string, error) {
	entries, err := os.ReadDir("projects")
	if err != nil {
		return nil, err
	}

	projects := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if projectNamePattern.MatchString(name) {
			projects = append(projects, name)
		}
	}
	sort.Strings(projects)
	return projects, nil
}

func ensureCommands(commands ...string) error {
	for _, name := range commands {
		if _, err := exec.LookPath(name); err != nil {
			return fmt.Errorf("%s 명령을 찾을 수 없습니다: %w", name, err)
		}
	}
	return nil
}

func resolveLocalBinaryPath(projectName string) (string, error) {
	localPath := filepath.Join(".", projectName)
	info, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("루트 디렉터리에 바이너리가 없습니다: %s", localPath)
		}
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("바이너리 경로가 디렉터리입니다: %s", localPath)
	}
	if err := validateLinuxELF(localPath, projectName); err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func validateLinuxELF(filePath, projectName string) error {
	// #nosec G304 -- 검증된 서비스명으로 만든 루트 경로(./{service})만 열어 확인한다.
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header := make([]byte, 4)
	if _, err := io.ReadFull(f, header); err != nil {
		return fmt.Errorf("바이너리 헤더를 읽지 못했습니다: %w", err)
	}

	if bytes.Equal(header, []byte{0x7f, 'E', 'L', 'F'}) {
		return nil
	}
	if bytes.Equal(header[:2], []byte{'M', 'Z'}) {
		return fmt.Errorf(
			"%s는 Windows 실행파일(PE)입니다. Linux 바이너리로 다시 빌드하세요. 예: GOOS=linux GOARCH=amd64 go build -o ./%s ./projects/%s/cmd",
			filePath, projectName, projectName,
		)
	}

	return fmt.Errorf("%s는 Linux ELF 실행파일이 아닙니다", filePath)
}

func writeServiceFile(servicePath, projectName, remoteAppDir string) error {
	content := fmt.Sprintf(`[Unit]
Description=%s Service
After=network.target

[Service]
ExecStart=%s
WorkingDirectory=%s
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
`, projectName, path.Join(remoteAppDir, projectName), remoteAppDir)

	return os.WriteFile(servicePath, []byte(content), 0600)
}

func renderDeployScript(projectName string, cfg DeployConfig, remoteBinaryTmp, remoteServiceTmp, remoteCaddyTmp string) string {
	remoteBinary := path.Join(cfg.RemoteAppDir, projectName)
	remoteService := path.Join(cfg.RemoteSystemdDir, projectName+".service")
	remoteSrvProject := path.Join(cfg.RemoteSrvRoot, projectName)
	remoteSrvData := path.Join(remoteSrvProject, "data")

	lines := []string{
		"set -euo pipefail",
		fmt.Sprintf("sudo mkdir -p %s", shellEscape(cfg.RemoteAppDir)),
		fmt.Sprintf("sudo install -m 755 %s %s", shellEscape(remoteBinaryTmp), shellEscape(remoteBinary)),
		fmt.Sprintf("sudo rm -f %s", shellEscape(remoteBinaryTmp)),
		fmt.Sprintf("sudo install -m 644 %s %s", shellEscape(remoteServiceTmp), shellEscape(remoteService)),
		fmt.Sprintf("sudo rm -f %s", shellEscape(remoteServiceTmp)),
		fmt.Sprintf("sudo install -m 644 %s %s", shellEscape(remoteCaddyTmp), shellEscape(cfg.RemoteCaddyfile)),
		fmt.Sprintf("sudo rm -f %s", shellEscape(remoteCaddyTmp)),
		"sudo systemctl daemon-reload",
		fmt.Sprintf("sudo mkdir -p %s", shellEscape(remoteSrvProject)),
		fmt.Sprintf("sudo mkdir -p %s", shellEscape(remoteSrvData)),
		fmt.Sprintf("sudo chown -R www-data:www-data %s", shellEscape(remoteSrvProject)),
		fmt.Sprintf("sudo chmod -R 755 %s", shellEscape(remoteSrvProject)),
		fmt.Sprintf("sudo systemctl enable %s", shellEscape(projectName+".service")),
		fmt.Sprintf("sudo systemctl restart %s", shellEscape(projectName+".service")),
		"sudo systemctl reload caddy",
		fmt.Sprintf("sudo systemctl --no-pager --full status %s || true", shellEscape(projectName+".service")),
	}
	return strings.Join(lines, "\n") + "\n"
}

func renderUndeployScript(projectName string, cfg DeployConfig, remoteCaddyTmp string, removeSrvData bool) string {
	remoteBinary := path.Join(cfg.RemoteAppDir, projectName)
	remoteService := path.Join(cfg.RemoteSystemdDir, projectName+".service")
	remoteSrvProject := path.Join(cfg.RemoteSrvRoot, projectName)

	lines := []string{
		"set -euo pipefail",
		fmt.Sprintf("sudo systemctl stop %s || true", shellEscape(projectName+".service")),
		fmt.Sprintf("sudo systemctl disable %s || true", shellEscape(projectName+".service")),
		fmt.Sprintf("sudo rm -f %s", shellEscape(remoteService)),
		"sudo systemctl daemon-reload",
		fmt.Sprintf("sudo systemctl reset-failed %s || true", shellEscape(projectName+".service")),
		fmt.Sprintf("sudo rm -f %s", shellEscape(remoteBinary)),
	}

	if removeSrvData {
		lines = append(lines, fmt.Sprintf("sudo rm -rf %s", shellEscape(remoteSrvProject)))
	}

	lines = append(lines,
		fmt.Sprintf("sudo install -m 644 %s %s", shellEscape(remoteCaddyTmp), shellEscape(cfg.RemoteCaddyfile)),
		fmt.Sprintf("sudo rm -f %s", shellEscape(remoteCaddyTmp)),
		"sudo systemctl reload caddy",
	)

	return strings.Join(lines, "\n") + "\n"
}

func runRemoteScript(cfg DeployConfig, scriptName, scriptBody string) error {
	localScriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("scaffold-%s-%d.sh", scriptName, time.Now().UnixNano()))
	fullScript := "#!/usr/bin/env bash\n" + scriptBody
	logInfo("원격 스크립트 파일 생성: %s", localScriptPath)
	if err := os.WriteFile(localScriptPath, []byte(fullScript), 0600); err != nil {
		return err
	}
	defer os.Remove(localScriptPath)

	remoteScriptPath := fmt.Sprintf("/tmp/scaffold-%s-%d.sh", scriptName, time.Now().UnixNano())
	logInfo("원격 스크립트 업로드 경로: %s", remoteScriptPath)
	if err := scpFile(cfg, localScriptPath, remoteScriptPath); err != nil {
		return err
	}
	defer func() {
		logInfo("원격 스크립트 정리: %s", remoteScriptPath)
		_ = runSSHCommand(cfg, fmt.Sprintf("rm -f %s", shellEscape(remoteScriptPath)))
	}()

	logInfo("원격 스크립트 실행 시작: %s", remoteScriptPath)
	return runSSHCommand(cfg, fmt.Sprintf("bash %s", shellEscape(remoteScriptPath)))
}

func scpFile(cfg DeployConfig, localPath, remotePath string) error {
	sshTarget := fmt.Sprintf("%s@%s:%s", cfg.SSHUser, cfg.SSHHost, remotePath)
	args := buildSCPCommonArgs(cfg)
	args = append(args, localPath, sshTarget)

	fmt.Printf("파일 전송: %s -> %s\n", localPath, sshTarget)
	return runCommand("scp", args...)
}

func runSSHCommand(cfg DeployConfig, remoteCommand string) error {
	args := buildSSHCommonArgs(cfg)
	args = append(args, fmt.Sprintf("%s@%s", cfg.SSHUser, cfg.SSHHost), remoteCommand)

	fmt.Printf("원격 명령 실행: %s\n", remoteCommand)
	return runCommand("ssh", args...)
}

func buildSSHCommonArgs(cfg DeployConfig) []string {
	args := []string{
		"-p",
		strconv.Itoa(cfg.SSHPort),
	}

	if cfg.DisableHostKeyTest {
		args = append(args, "-o", "StrictHostKeyChecking=no")
	}
	if !cfg.SkipSSHKey {
		args = append(args, "-i", cfg.SSHKeyPath)
	}
	return args
}

func buildSCPCommonArgs(cfg DeployConfig) []string {
	args := []string{
		"-P",
		strconv.Itoa(cfg.SSHPort),
	}

	if cfg.DisableHostKeyTest {
		args = append(args, "-o", "StrictHostKeyChecking=no")
	}
	if !cfg.SkipSSHKey {
		args = append(args, "-i", cfg.SSHKeyPath)
	}
	return args
}

func runCommand(name string, args ...string) error {
	logInfo("명령 실행: %s", formatCommandForLog(name, args))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s 명령 실행 실패: %w", name, err)
	}
	logInfo("명령 완료: %s", name)
	return nil
}

func projectNameValidator(answer interface{}) error {
	s, ok := answer.(string)
	if !ok {
		return fmt.Errorf("문자열만 입력할 수 있습니다")
	}
	s = strings.TrimSpace(s)
	if !projectNamePattern.MatchString(s) {
		return fmt.Errorf("프로젝트 명은 소문자/숫자/하이픈만 사용할 수 있습니다")
	}
	return nil
}

func portValidator(answer interface{}) error {
	raw, ok := answer.(string)
	if !ok {
		return fmt.Errorf("숫자 포트만 입력할 수 있습니다")
	}
	portRaw := strings.TrimSpace(raw)
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return fmt.Errorf("포트는 숫자로 입력해야 합니다")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("포트는 1~65535 범위여야 합니다")
	}
	return nil
}

func shellEscape(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

func logTitle(message string) {
	fmt.Printf("\n[작업] %s\n", message)
}

func logStep(message string) {
	fmt.Printf("[단계] %s\n", message)
}

func logSuccess(message string) {
	fmt.Printf("[완료] %s\n", message)
}

func logInfo(format string, args ...interface{}) {
	fmt.Printf("[정보 %s] %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

func formatCommandForLog(name string, args []string) string {
	items := make([]string, 0, len(args)+1)
	items = append(items, name)

	maskNext := false
	for _, arg := range args {
		if maskNext {
			items = append(items, "***")
			maskNext = false
			continue
		}
		if arg == "-i" {
			items = append(items, arg)
			maskNext = true
			continue
		}
		items = append(items, arg)
	}

	return strings.Join(items, " ")
}
