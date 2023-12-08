
document.addEventListener('DOMContentLoaded', (event) => {
  let params = getParams();
  disableElements(params.cb);
  setFocus(params.cb);
});

function setFocus(eid) {
  // var audio = new Audio('https://media.geeksforgeeks.org/wp-content/uploads/20190531135120/beep.mp3');
  let tb;
  if (eid === "lhs-button1") {
    // audio.play();
    tb = document.getElementById("fd1-n");
  } else if (eid === "lhs-button2") {
    tb = document.getElementById("fd2-n");
  }
  tb.focus();
  /***
  input type="number" doesn't support setSelectionRange.
  You can use type="text" and inputmode="numeric". This will show a numeric keyboard for
  mobile users and supports setSelectionRange.
  ***/
  // tb.type = "text";
  tb.setSelectionRange(0, tb.value.length);
  // tb.type = "number";
  // console.log(`Position start: ${tb.selectionStart}`);
  // console.log(`Position end: ${tb.selectionEnd}`);
}

function disableElements(eid) {
  if (eid === "lhs-button1") {
    document.getElementById("lhs-button1").disabled = true;
    document.getElementById("lhs-button2").disabled = true;
  } else if (eid === "lhs-button2") {
    document.getElementById("lhs-button1").disabled = true;
    document.getElementById("lhs-button2").disabled = true;
  }
}
