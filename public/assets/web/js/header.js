const mobileMenuToggler = document.getElementById('mobile-menu-toggler');
const mobileMenu = document.getElementById('mobile-menu');

mobileMenuToggler.addEventListener('click', () => {
  document.body.classList.toggle('dee--within-screen');
  mobileMenuToggler.classList.toggle('dee--open');
  mobileMenu.classList.toggle('dee--open');
});
