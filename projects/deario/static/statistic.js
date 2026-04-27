;(function () {
  const moodSets = [
    { label: "😁", key: "mood1", color: "#ffeb3b" },
    { label: "🙂", key: "mood2", color: "#8bc34a" },
    { label: "😐", key: "mood3", color: "#03a9f4" },
    { label: "😣", key: "mood4", color: "#ff9800" },
    { label: "😭", key: "mood5", color: "#f44336" },
  ]

  document.addEventListener("DOMContentLoaded", initStatisticPage)

  async function initStatisticPage() {
    const countChart = document.getElementById("countChart")
    const moodChart = document.getElementById("moodStackChart")
    if (!countChart || !moodChart) return

    try {
      const data = await fetchStatisticData()
      const labels = data.months.map(formatMonth)

      renderDiaryCountChart(countChart, labels, data)
      renderMoodStackChart(moodChart, labels, data)
    } catch (err) {
      console.error(err)
    }
  }

  async function fetchStatisticData() {
    const response = await fetch("/statistic/data")
    if (!response.ok) {
      throw new Error("통계 데이터를 불러오지 못했습니다.")
    }
    return response.json()
  }

  function renderDiaryCountChart(canvas, labels, data) {
    new Chart(canvas.getContext("2d"), {
      type: "bar",
      data: {
        labels,
        datasets: [
          {
            label: "작성 수",
            data: data.diaryCount,
            backgroundColor: "rgba(33,150,243,0.5)",
          },
        ],
      },
    })
  }

  function renderMoodStackChart(canvas, labels, data) {
    new Chart(canvas.getContext("2d"), {
      type: "bar",
      data: {
        labels,
        datasets: moodSets.map((mood) => ({
          label: mood.label,
          data: data[mood.key],
          backgroundColor: mood.color,
        })),
      },
      options: {
        scales: {
          x: { stacked: true },
          y: { stacked: true, beginAtZero: true },
        },
      },
    })
  }

  function formatMonth(monthStr) {
    const year = monthStr.substring(2, 4)
    const month = parseInt(monthStr.substring(4, 6), 10)
    return `${year}년 ${month}월`
  }
})()
