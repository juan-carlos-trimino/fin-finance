
document.addEventListener('DOMContentLoaded', (event) => {
  let params = getParams();
  disableElements(params.cb);
  setFocus(params.cb);
});

function setFocus(eid) {
  let tb;
  if (eid === "lhs-button1") {
    tb = document.getElementById("fd1-n");
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
  }
}
