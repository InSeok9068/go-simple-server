document.addEventListener("DOMContentLoaded", () => {
  const el = document.getElementById("calendar-picker");
  if (!el) return;
  const fp = flatpickr(el, {
    inline: true,
    dateFormat: "Ymd",
    onChange: function (_sel, dateStr) {
      if (dateStr) {
        location.href = "/?date=" + dateStr;
      }
    },
    onMonthChange: function (_sel, _dateStr, instance) {
      loadDiaryDates(instance);
    },
    onReady: function (_sel, _dateStr, instance) {
      loadDiaryDates(instance);
    },
  });

  async function loadDiaryDates(instance) {
    const year = instance.currentYear;
    const month = String(instance.currentMonth + 1).padStart(2, "0");
    try {
      const res = await fetch(`/diary/month?month=${year}${month}`);
      if (!res.ok) return;
      const dates = await res.json();
      highlightDates(instance, dates);
    } catch (e) {
      console.error(e);
    }
  }

  function highlightDates(instance, dates) {
    instance.days
      .querySelectorAll(".has-diary")
      .forEach((e) => e.classList.remove("has-diary"));
    dates.forEach((d) => {
      const selector = `[aria-label="${d.slice(0, 4)}-${d.slice(4, 6)}-${d.slice(6, 8)}"]`;
      const dayElem = instance.days.querySelector(selector);
      if (dayElem) {
        dayElem.classList.add("has-diary");
      }
    });
  }
});
