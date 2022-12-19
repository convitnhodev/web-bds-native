const sendCode = document.getElementById('send-code');
const countdownTimer = document.getElementById('countdown-timer');
const verifyForm = document.getElementById('verify-form');

function getSeconds(countdown) {
  const distance = countdown - new Date().getTime();

  return Math.floor((distance / 1000) % 1000);
}

async function requestCode() {
  try {
    const token = localStorage.getItem('token');
    const resp = await fetch('/api/v1/accounts/verification', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });

    if (resp.ok) {
      return;
    }

    throw new Error();
  } catch (error) {
    console.log(error);
  }
}

sendCode.addEventListener('click', () => {
  sendCode.setAttribute('disabled', '');
  requestCode();

  const countdown = new Date().getTime() + 60000 + 200; // padding 200ms

  countdownTimer.innerHTML = `(${getSeconds(countdown)})`;

  const interval = setInterval(() => {
    const seconds = getSeconds(countdown);

    countdownTimer.innerHTML = `(${seconds})`;

    if (seconds < 0) {
      clearInterval(interval);
      sendCode.removeAttribute('disabled');

      countdownTimer.innerHTML = '';
    }
  }, 1000);
});

verifyForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const verifyFormOverlay = document.getElementById('verify-form-overlay');
  const verifyFormSubmit = document.getElementById('verify-form-submit');
  const verifyFormErr = document.getElementById('verify-form-error');

  // reset
  verifyFormErr.classList.add('hidden');

  const codeEl = document.getElementById('code');

  let code;

  if (codeEl instanceof HTMLInputElement) {
    code = codeEl.value;
  }

  verifyFormOverlay.classList.remove('hidden');
  verifyFormSubmit.setAttribute('disabled', '');

  // submit
  try {
    const token = localStorage.getItem('token');
    const resp = await fetch(`/api/v1/accounts/verification/${code}/code`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });

    if (resp.ok) {
      const json = await resp.json();

      localStorage.setItem('token', json.token);

      location.href = '/';

      return;
    }

    throw new Error();
  } catch (error) {
    verifyFormOverlay.classList.add('hidden');
    verifyFormSubmit.removeAttribute('disabled');
    verifyFormErr.classList.remove('hidden');
  }
});
