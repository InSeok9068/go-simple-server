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
    // 서버에서 받은 일기 날짜를 집합(set)으로 만들어, 일기 날짜 여부를 빠르게 확인할 수 있도록.
    const diaryDatesSet = new Set(dates);

    // 현재 달력 뷰에 표시되는 각 날짜 요소를 순회.
    // flatpickr-day 클래스는 선택 가능한 날짜 요소.
    instance.days.querySelectorAll(".flatpickr-day").forEach((dayElem) => {
      // 요소에 연관된 날짜 객체를 가져옴.
      const dateObj = dayElem.dateObj;
      if (!dateObj) return;

      // 날짜 객체를 'YYYYMMDD' 문자열로 포맷하여, 서버에서 받은 날짜와 일치하는지 확인.
      const year = dateObj.getFullYear();
      const month = String(dateObj.getMonth() + 1).padStart(2, "0");
      const day = String(dayElem.dateObj.getDate()).padStart(2, "0");
      const dateStr = `${year}${month}${day}`;

      // 일기 날짜 집합에 포함된 경우 'has-diary' 클래스를 추가,
      // 포함되지 않은 경우 'has-diary' 클래스를 제거.
      if (diaryDatesSet.has(dateStr)) {
        dayElem.classList.add("has-diary");
      } else {
        dayElem.classList.remove("has-diary");
      }
    });
  }
});
