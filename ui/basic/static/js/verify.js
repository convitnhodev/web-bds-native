const sendCode = document.getElementById('send-code');
const countdownTimer = document.getElementById('countdown-timer');
const valueInput = document.getElementById('verify-value');
const errorMsg = document.getElementById('verify-error');

const EMAIL_RX =
  /^[a-z0-9!#$%&'*+\/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+\/=?^_{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/;
const PHONE_RX = /((^(\+84|84|0|0084){1})(3|5|7|8|9))+([0-9]{8})$/;

function getSeconds(countdown) {
  const distance = countdown - new Date().getTime();

  return Math.floor((distance / 1000) % 1000);
}

async function requestCode(type, value) {
  try {
    const resp = await fetch(`/ajax/verify-${type}?${type}=${value}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
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

sendCode.addEventListener('click', async () => {
  if (!(valueInput instanceof HTMLInputElement)) {
    return;
  }

  const value = valueInput.value;

  if (!value) {
    errorMsg.classList.remove('hidden');

    return;
  }

  const type = sendCode.getAttribute('data-verify');

  let valid = false;

  if (type === 'email' && EMAIL_RX.test(value)) {
    valid = true;
  }

  if (type === 'phone' && PHONE_RX.test(value)) {
    valid = true;
  }

  if (!valid) {
    errorMsg.classList.remove('hidden');

    return;
  }

  sendCode.setAttribute('disabled', '');

  await requestCode(type, value).catch(() => {
    // do nothing
  });

  errorMsg.classList.add('hidden');

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
