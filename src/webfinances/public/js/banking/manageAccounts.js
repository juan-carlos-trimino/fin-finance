
document.addEventListener('DOMContentLoaded', (event) => {
  let params = getParams();
  disableElements(params.cb);
  setFocus(params.cb);
});

function setFocus(eid) {
  // var audio = new Audio('https://media.geeksforgeeks.org/wp-content/uploads/20190531135120/beep.mp3');
  let tb = null;
  if (eid === "lhs-button1") {
    // audio.play();
    tb = document.getElementById("fd1-bankname");
  }
  //
  if (tb !== null) {
    tb.focus();
    /***
    input type="number" doesn't support setSelectionRange.
    You can use type="text" and inputmode="numeric". This will show a numeric keyboard for
    mobile users and supports setSelectionRange.
    ***/
    // tb.type = "text";
    tb.setSelectionRange(0, tb.value.length);
  }
}

function disableElements(eid) {
  // var audio = new Audio('https://media.geeksforgeeks.org/wp-content/uploads/20190531135120/beep.mp3');
  if (eid === "lhs-button1") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = true;
    document.getElementById("lhs-button2").disabled = false;
  } else if (eid === "lhs-button2") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = true;
  }
}
