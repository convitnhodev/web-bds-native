const disclosures = document.querySelectorAll('[data-disclosure="item"]');

disclosures.forEach((el) => {
  const toggler = el.querySelector('[data-disclosure="toggler"]');
  const content = el.querySelector('[data-disclosure="content"]');

  if (toggler && toggler instanceof HTMLElement) {
    toggler.addEventListener('click', () => {
      toggler.classList.toggle('dee--open');

      if (content && content instanceof HTMLElement) {
        content.classList.toggle('dee--open');
      }
    });
  }
});
