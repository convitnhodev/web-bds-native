function getInt(attr, el) {
  return parseInt(el.getAttribute(attr), 10) || 0;
}

function getFloat(attr, el) {
  return parseFloat(el.getAttribute(attr)) || 0;
}

function getInputAsInt(el) {
  if (el instanceof HTMLInputElement) {
    return parseInt(el.value, 10) || 0;
  }

  return 0;
}

const formatter = new Intl.NumberFormat('vi-VN', {
  currency: 'VND',
  currencyDisplay: 'code',
  style: 'currency',
});

const slotRemainEl = document.getElementById('slot-remain');
const slotQtyEl = document.getElementById('slot-qty');
const slotPriceEl = document.getElementById('slot-price');
const slotSubEl = document.getElementById('slot-sub');
const slotAddEl = document.getElementById('slot-add');
const totalEl = document.getElementById('total');
const depositEl = document.getElementById('deposit');

const slotRemainValue = getInt('data-slot-remain', slotRemainEl);
const slotPriceValue = getInt('data-slot-price', slotPriceEl);
const depositValue = getFloat('data-deposit', depositEl);

function calcTotal(_qty) {
  let qty = _qty;

  if (!qty) {
    qty = getInputAsInt(slotQtyEl);
  }

  return qty * slotPriceValue;
}

function updateDeposit(qty) {
  const total = calcTotal(qty);

  depositEl.innerHTML = formatter.format((total * depositValue) / 100);
}

function updateTotal(qty) {
  const total = calcTotal(qty);

  totalEl.innerHTML = formatter.format(total);
}

function updateSlotPrice() {
  slotPriceEl.innerHTML = formatter.format(slotPriceValue);
}

document.addEventListener('DOMContentLoaded', () => {
  updateTotal();
  updateDeposit();
  updateSlotPrice();
});

slotQtyEl.addEventListener('change', () => {
  updateTotal();
  updateDeposit();
});

slotSubEl.addEventListener('click', () => {
  const qty = getInputAsInt(slotQtyEl);

  if (qty === 0) {
    return;
  }

  const nqty = qty - 1;

  slotQtyEl.value = nqty;

  updateTotal(nqty);
  updateDeposit(nqty);
});

slotAddEl.addEventListener('click', () => {
  const qty = getInputAsInt(slotQtyEl);

  if (qty === slotRemainValue) {
    return;
  }

  const nqty = qty + 1;

  slotQtyEl.value = nqty;

  updateTotal(nqty);
  updateDeposit(nqty);
});
