const loginForm = document.getElementById('login-form');

loginForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const loginFormOverlay = document.getElementById('login-form-overlay');
  const loginFormSubmit = document.getElementById('login-form-submit');
  const loginFormErr = document.getElementById('login-form-error');
  const phoneNumberErr = document.getElementById('phone-number-error');

  // reset
  loginFormErr.classList.add('hidden');
  phoneNumberErr.classList.add('hidden');

  const phoneNumberEl = document.getElementById('phone-number');
  const passwordEl = document.getElementById('password');

  let phoneNumber;
  let password;

  if (phoneNumberEl instanceof HTMLInputElement) {
    phoneNumber = phoneNumberEl.value;
  }

  if (passwordEl instanceof HTMLInputElement) {
    password = passwordEl.value;
  }

  // validate
  let err = false;

  if (!isPhoneNumber(phoneNumber)) {
    err = true;

    phoneNumberErr.classList.remove('hidden');
  }

  if (err) {
    return;
  }

  loginFormOverlay.classList.remove('hidden');
  loginFormSubmit.setAttribute('disabled', '');

  const body = {
    password: password,
    phone_number: phoneNumber.startsWith('0')
      ? phoneNumber.replace('0', '+84')
      : phoneNumber,
  };

  // submit
  try {
    const resp = await fetch('/api/v1/accounts/auth', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    let token;

    if (resp.ok) {
      const json = await resp.json();

      console.log(json);

      localStorage.setItem('token', json.token);

      token = json.token;
    }

    if (!token) {
      throw new Error();
    }

    const profileResp = await fetch('/api/v1/accounts/profile', {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });

    if (profileResp.ok) {
      const json = await profileResp.json();

      if (json.is_verified) {
        location.href = '/';
      } else {
        location.href = '/verify';
      }

      return;
    }

    throw new Error();
  } catch (error) {
    console.log(error);
    loginFormOverlay.classList.add('hidden');
    loginFormSubmit.removeAttribute('disabled');
    loginFormErr.classList.remove('hidden');
  }
});
