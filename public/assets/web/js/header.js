const mobileMenuToggler = document.getElementById('mobile-menu-toggler');
const mobileMenu = document.getElementById('mobile-menu');

mobileMenuToggler.addEventListener('click', () => {
  document.body.classList.toggle('dee--within-screen');
  mobileMenuToggler.classList.toggle('dee--open');
  mobileMenu.classList.toggle('dee--open');
});

const mdMm = window.matchMedia('(min-width: 768px)');

mdMm.addEventListener('change', (e) => {
  if (e.matches) {
    document.body.classList.remove('dee--within-screen');
    mobileMenuToggler.classList.remove('dee--open');
    mobileMenu.classList.remove('dee--open');
  }
});
