
document.addEventListener('DOMContentLoaded', (event) => {
  let params = getParams();
  disableElements(params.cb);
  setFocus(params.cb);
});

function setFocus(eid) {
  let tb;
  if (eid === "lhs-button1") {
    tb = document.getElementById("fd1-nominal");
  } else if (eid === "lhs-button2") {
    tb = document.getElementById("fd2-effective");
  } else if (eid === "lhs-button3") {
    tb = document.getElementById("fd3-nominal");
  } else if (eid === "lhs-button4") {
    tb = document.getElementById("fd4-interest");
  } else if (eid === "lhs-button5") {
    tb = document.getElementById("fd5-values");
  } else if (eid === "lhs-button6") {
    tb = document.getElementById("fd6-time");
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
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
  } else if (eid === "lhs-button2") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = true;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
  } else if (eid === "lhs-button3") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = true;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
  } else if (eid === "lhs-button4") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = true;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
  } else if (eid === "lhs-button5") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
  } else if (eid === "lhs-button6") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = true;
  }
}
