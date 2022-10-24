const registerForm = document.getElementById('register-form');

registerForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const registerFormOverlay = document.getElementById('register-form-overlay');
  const registerFormSubmit = document.getElementById('register-form-submit');
  const registerFormErr = document.getElementById('register-form-error');
  const phoneNumberErr = document.getElementById('phone-number-error');
  const passwordConfirmationErr = document.getElementById(
    'password-confirmation-error'
  );

  // reset
  registerFormErr.classList.add('hidden');
  phoneNumberErr.classList.add('hidden');
  passwordConfirmationErr.classList.add('hidden');

  const firstNameEl = document.getElementById('first-name');
  const lastNameEl = document.getElementById('last-name');
  const phoneNumberEl = document.getElementById('phone-number');
  const passwordEl = document.getElementById('password');
  const passwordConfirmationEl = document.getElementById(
    'password-confirmation'
  );

  let firstName;
  let lastName;
  let phoneNumber;
  let password;
  let passwordConfirmation;

  if (firstNameEl instanceof HTMLInputElement) {
    firstName = firstNameEl.value;
  }

  if (lastNameEl instanceof HTMLInputElement) {
    lastName = lastNameEl.value;
  }

  if (phoneNumberEl instanceof HTMLInputElement) {
    phoneNumber = phoneNumberEl.value;
  }

  if (passwordEl instanceof HTMLInputElement) {
    password = passwordEl.value;
  }

  if (passwordConfirmationEl instanceof HTMLInputElement) {
    passwordConfirmation = passwordConfirmationEl.value;
  }

  // validate
  let err = false;

  if (!isPhoneNumber(phoneNumber)) {
    err = true;

    phoneNumberErr.classList.remove('hidden');
  }

  if (password !== passwordConfirmation) {
    err = true;

    passwordConfirmationErr.classList.remove('hidden');
  }

  if (err) {
    return;
  }

  registerFormOverlay.classList.remove('hidden');
  registerFormSubmit.setAttribute('disabled', '');

  const body = {
    first_name: firstName,
    last_name: lastName,
    password: password,
    password_confirmation: passwordConfirmation,
    phone_number: phoneNumber.startsWith('0')
      ? phoneNumber.replace('0', '+84')
      : phoneNumber,
  };

  // submit
  try {
    const resp = await fetch('/api/v1/accounts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    if (resp.ok) {
      await resp.json();

      location.href = '/login';

      return;
    }

    throw new Error();
  } catch (error) {
    registerFormErr.classList.remove('hidden');
    registerFormOverlay.classList.add('hidden');
    registerFormSubmit.removeAttribute('disabled');
  }
});
