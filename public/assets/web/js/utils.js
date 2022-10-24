function isPhoneNumber(s) {
  return /((^(\+84|84|0|0084){1})(3|5|7|8|9))+([0-9]{8})$/.test(s);
}
