;(function () {
  document.addEventListener("DOMContentLoaded", initCalendar)

  function initCalendar() {
    const el = document.getElementById("calendar-picker")
    if (!el) return

    flatpickr(el, {
      inline: true,
      dateFormat: "Ymd",
      appendTo: el,
      locale: "ko",
      defaultDate: selectedDate(),
      onChange: (_selectedDates, dateStr) => {
        if (dateStr) {
          location.href = `/?date=${dateStr}`
        }
      },
      onMonthChange: (_selectedDates, _dateStr, instance) => {
        loadDiaryDates(instance)
      },
      onReady: (_selectedDates, _dateStr, instance) => {
        loadDiaryDates(instance)
      },
    })
  }

  function selectedDate() {
    return new URLSearchParams(location.search).get("date") || "today"
  }

  async function loadDiaryDates(instance) {
    const year = instance.currentYear
    const month = String(instance.currentMonth + 1).padStart(2, "0")

    try {
      const response = await fetch(`/diary/month?month=${year}${month}`)
      if (!response.ok) return

      highlightDates(instance, await response.json())
    } catch (err) {
      console.error(err)
    }
  }

  function highlightDates(instance, dates) {
    const diaryDates = new Set(dates)

    instance.days.querySelectorAll(".flatpickr-day").forEach((dayEl) => {
      const date = dayEl.dateObj
      if (!date) return

      dayEl.classList.toggle("has-diary", diaryDates.has(formatDate(date)))
    })
  }

  function formatDate(date) {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, "0")
    const day = String(date.getDate()).padStart(2, "0")
    return `${year}${month}${day}`
  }
})()
