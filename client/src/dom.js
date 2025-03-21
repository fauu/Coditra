export function scrollTo(el, { topOffset, smooth }) {
  const pos = window.pageYOffset + el.getBoundingClientRect().top;
  const offsetPos = pos - topOffset;

  window.scrollTo({
    top: offsetPos,
    behavior: smooth ? "smooth" : "auto",
  });
}
