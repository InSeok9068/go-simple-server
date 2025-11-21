console.log("statistic page loaded");

const moodIcons = ["", "😁", "🙂", "😐", "😣", "😭"];

document.addEventListener("DOMContentLoaded", () => {
  fetch("/statistic/data")
    .then((res) => res.json())
    .then((data) => {
      const formattedMonths = data.months.map((monthStr) => {
        const year = monthStr.substring(2, 4);
        const month = parseInt(monthStr.substring(4, 6), 10);
        return `${year}년 ${month}월`;
      });

      const countCtx = document.getElementById("countChart").getContext("2d");
      new Chart(countCtx, {
        type: "bar",
        data: {
          labels: formattedMonths,
          datasets: [
            {
              label: "작성 수",
              data: data.diaryCount,
              backgroundColor: "rgba(33,150,243,0.5)",
            },
          ],
        },
      });

      const stackCtx = document
        .getElementById("moodStackChart")
        .getContext("2d");
      new Chart(stackCtx, {
        type: "bar",
        data: {
          labels: formattedMonths,
          datasets: [
            {
              label: `${moodIcons[1]}`,
              data: data.mood1,
              backgroundColor: "#ffeb3b",
            },
            {
              label: `${moodIcons[2]}`,
              data: data.mood2,
              backgroundColor: "#8bc34a",
            },
            {
              label: `${moodIcons[3]}`,
              data: data.mood3,
              backgroundColor: "#03a9f4",
            },
            {
              label: `${moodIcons[4]}`,
              data: data.mood4,
              backgroundColor: "#ff9800",
            },
            {
              label: `${moodIcons[5]}`,
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
