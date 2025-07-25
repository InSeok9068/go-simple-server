#!/bin/bash

# sqlc 코드를 생성합니다.
# 인자로 특정 프로젝트를 지정하거나, 지정하지 않으면 모든 프로젝트에 대해 실행합니다.

PROJECT_NAME=$1
PROJECTS_DIR="projects"
SQLC_FILE="sqlc.yaml"

# 색상 정의
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 특정 프로젝트에 대해 실행
if [ -n "$PROJECT_NAME" ]; then
    PROJECT_PATH="$PROJECTS_DIR/$PROJECT_NAME"
    SQLC_PATH="$PROJECT_PATH/$SQLC_FILE"

    if [ -f "$SQLC_PATH" ]; then
        echo -e "${BLUE}Generating SQLC for project: ${GREEN}$PROJECT_NAME${NC}"
        sqlc generate -f "$SQLC_PATH"
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Successfully generated SQLC for ${PROJECT_NAME}${NC}"
        else
            echo -e "${RED}Failed to generate SQLC for ${PROJECT_NAME}${NC}"
        fi
    else
        echo -e "${RED}Error: ${YELLOW}$SQLC_PATH${RED} not found.${NC}"
        exit 1
    fi
else
    # 모든 프로젝트에 대해 실행
    echo -e "${BLUE}Generating SQLC for all projects...${NC}"
    
    found_any=false
    for project in "$PROJECTS_DIR"/*; do
        if [ -d "$project" ]; then
            PROJECT_PATH="$project"
            SQLC_PATH="$PROJECT_PATH/$SQLC_FILE"

            if [ -f "$SQLC_PATH" ]; then
                found_any=true
                PROJECT_NAME=$(basename "$project")
                echo -e "\n${BLUE}Generating SQLC for project: ${GREEN}$PROJECT_NAME${NC}"
                sqlc generate -f "$SQLC_PATH"
                if [ $? -eq 0 ]; then
                    echo -e "${GREEN}Successfully generated SQLC for ${PROJECT_NAME}${NC}"
                else
                    echo -e "${RED}Failed to generate SQLC for ${PROJECT_NAME}${NC}"
                fi
            fi
        fi
    done

    if [ "$found_any" = false ]; then
        echo -e "${YELLOW}No projects with a sqlc.yaml file found.${NC}"
    fi
fi

echo -e "\n${GREEN}SQLC generation process finished.${NC}"
