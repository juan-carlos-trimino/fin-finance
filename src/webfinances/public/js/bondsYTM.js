
function onchangeRadio(radio) {
  if (radio === 'fd4-curinterest') {
    document.getElementById("fd4-curinterest").checked = true;
  } else {
    document.getElementById("fd4-bondprice").checked = true;
  }
  return;
}
