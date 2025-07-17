console.log("statistic page loaded");

const moodIcons = ["", "😁", "🙂", "😐", "😣", "😭"];

document.addEventListener("DOMContentLoaded", () => {
  fetch("/statistic/data")
    .then((res) => res.json())
    .then((data) => {
      const countCtx = document.getElementById("countChart").getContext("2d");
      new Chart(countCtx, {
        type: "bar",
        data: {
          labels: data.months,
          datasets: [
            {
              label: "작성 수",
              data: data.diaryCount,
              backgroundColor: "rgba(33,150,243,0.5)",
            },
          ],
        },
      });

      const moodCtx = document.getElementById("moodChart").getContext("2d");
      new Chart(moodCtx, {
        type: "line",
        data: {
          labels: data.months,
          datasets: [
            {
              label: "평균 기분",
              data: data.moodAvg,
              borderColor: "rgba(244,67,54,0.8)",
              fill: false,
            },
          ],
        },
        options: {
          scales: { y: { suggestedMin: 1, suggestedMax: 5 } },
        },
      });

      const stackCtx = document
        .getElementById("moodStackChart")
        .getContext("2d");
      new Chart(stackCtx, {
        type: "bar",
        data: {
          labels: data.months,
          datasets: [
            {
              label: `1 ${moodIcons[1]}`,
              data: data.mood1,
              backgroundColor: "#ffeb3b",
            },
            {
              label: `2 ${moodIcons[2]}`,
              data: data.mood2,
              backgroundColor: "#8bc34a",
            },
            {
              label: `3 ${moodIcons[3]}`,
              data: data.mood3,
              backgroundColor: "#03a9f4",
            },
            {
              label: `4 ${moodIcons[4]}`,
              data: data.mood4,
              backgroundColor: "#ff9800",
            },
            {
              label: `5 ${moodIcons[5]}`,
              data: data.mood5,
              backgroundColor: "#f44336",
            },
          ],
        },
        options: {
          scales: {
            x: { stacked: true },
            y: { stacked: true, beginAtZero: true },
          },
        },
      });
    })
    .catch((err) => console.error(err));
});
