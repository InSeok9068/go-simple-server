window.applyThemeColor = function (color) {
  try {
    const theme = materialDynamicColors.themeFromSourceColor(color);
    applyTheme(theme);
  } catch (e) {
    console.error("set color error", e);
  }
};
