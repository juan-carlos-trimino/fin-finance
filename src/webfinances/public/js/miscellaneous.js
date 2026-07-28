
document.addEventListener('DOMContentLoaded', (event) => {
  console.log("Entering document.addEventListener...");
  let params = getParams();
  disableElements(params.cb);
  setFocus(params.cb);
  console.log("Exiting document.addEventListener...");
});

function setFocus(eid) {
  console.log("Entering setFocus...");
  let tb = null;
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
  } else if (eid === "lhs-button7") {
    tb = document.getElementById("fd7-time");
  }
  //Only call .focus() if the element was actually found.
  if (tb) {
    tb.focus();
    /***
    input type="number" doesn't support setSelectionRange.
    You can use type="text" and inputmode="numeric". This will show a numeric keyboard for mobile users and supports setSelectionRange.
    Add a type check: Wrap the JavaScript code in a condition so it only runs on valid input types.
    ***/
    if (tb.type !== 'number' && typeof tb.setSelectionRange === 'function') {
      tb.setSelectionRange(0, tb.value.length);
    }
  } else {
    console.warn("No matching element found for ID:", eid);
  }
  console.log("Exiting setFocus...");
}

function disableElements(eid) {
  console.log("Entering disableElements...");
  if (eid === "lhs-button1") {
    document.getElementById("lhs-button1").disabled = true;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
  } else if (eid === "lhs-button2") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = true;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
  } else if (eid === "lhs-button3") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = true;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
  } else if (eid === "lhs-button4") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = true;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
  } else if (eid === "lhs-button5") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
  } else if (eid === "lhs-button6") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = true;
    document.getElementById("lhs-button7").disabled = false;
  } else if (eid === "lhs-button7") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = false;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = true;
  }
  console.log("Exiting disableElements...");
}
