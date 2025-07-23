window.applyThemeColor = async (color) => {
  try {
    await window.ui("theme", color);
  } catch (e) {
    console.error("set color error", e);
  }
};
