
document.addEventListener('DOMContentLoaded', (event) => {
  let params = getParams();
  disableElements(params.cb);
  setFocus(params.cb);
  amountOfInterest(params.leap);
});

function setFocus(eid) {
  let tb;
  if (eid === "lhs-button1") {
    tb = document.getElementById("fd1-time");
  } else if (eid === "lhs-button2") {
    tb = document.getElementById("fd2-time");
  } else if (eid === "lhs-button3") {
    tb = document.getElementById("fd3-time");
  } else if (eid === "lhs-button4") {
    tb = document.getElementById("fd4-interest");
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
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
  } else if (eid === "lhs-button2") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = true;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
  } else if (eid === "lhs-button3") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = true;
    document.getElementById("lhs-button4").disabled = false;
  } else if (eid === "lhs-button4") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = true;
  }
}

function amountOfInterest(leap) {
  if (leap === "on") {
    document.getElementById("fd1-leap").checked = true;
  } else {
    document.getElementById("fd1-leap").checked = false;
  }
}
