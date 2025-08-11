document.addEventListener("DOMContentLoaded", () => {
  const el = document.getElementById("calendar-picker");
  if (!el) return;

  // URL에서 'date' 파라미터를 읽어와 달력의 초기 날짜로 설정합니다.
  const urlParams = new URLSearchParams(window.location.search);
  const dateParam = urlParams.get("date");

  const fp = flatpickr(el, {
    // --- 기본 설정 ---
    inline: true, // 달력을 항상 표시합니다.
    dateFormat: "Ymd", // 날짜 형식 (예: 20231231)
    appendTo: el, // 달력이 삽입될 부모 요소
    locale: "ko", // 한국어 지원

    // --- UX 개선 옵션 ---
    // URL에 날짜 파라미터가 있으면 해당 날짜를, 없으면 오늘 날짜를 기본으로 보여줍니다.
    defaultDate: dateParam || "today",

    // --- 이벤트 핸들러 ---
    onChange: function (_sel, dateStr) {
      // 날짜를 선택하면 해당 날짜의 일기 페이지로 이동합니다.
      if (dateStr) {
        location.href = "/?date=" + dateStr;
      }
    },
    onMonthChange: function (_sel, _dateStr, instance) {
      // 월이 변경될 때마다 해당 월의 일기 작성 기록을 불러옵니다.
      loadDiaryDates(instance);
    },
    onReady: function (_sel, _dateStr, instance) {
      // 달력이 준비되면 현재 월의 일기 작성 기록을 불러옵니다.
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
