document.addEventListener("DOMContentLoaded", () => {
  const el = document.getElementById("calendar-picker");
  if (!el) return;
  const fp = flatpickr(el, {
    inline: true,
    dateFormat: "Ymd",
    appendTo: el,
    locale: "ko",
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
    // Create a Set of date strings from the server for efficient lookup.
    const diaryDatesSet = new Set(dates);

    // Iterate over each day element in the current calendar view.
    // flatpickr-day is the class for selectable day elements.
    instance.days.querySelectorAll(".flatpickr-day").forEach((dayElem) => {
      // Get the JavaScript Date object associated with the element.
      const dateObj = dayElem.dateObj;
      if (!dateObj) return;

      // Format the date object into 'YYYYMMDD' string to match the server's format.
      const year = dateObj.getFullYear();
      const month = String(dateObj.getMonth() + 1).padStart(2, "0");
      const day = String(dayElem.dateObj.getDate()).padStart(2, "0");
      const dateStr = `${year}${month}${day}`;

      // Add or remove the 'has-diary' class based on whether the date is in our set.
      if (diaryDatesSet.has(dateStr)) {
        dayElem.classList.add("has-diary");
      } else {
        dayElem.classList.remove("has-diary");
      }
    });
  }
});
