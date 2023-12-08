
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
    tb = document.getElementById("fd1-taxfree");
  } else if (eid === "lhs-button2") {
    tb = document.getElementById("fd2-facevalue");
  } else if (eid === "lhs-button3") {
    tb = document.getElementById("fd3-facevalue");
  } else if (eid === "lhs-button4") {
    tb = document.getElementById("fd4-facevalue");
  } else if (eid === "lhs-button5") {
    tb = document.getElementById("fd5-facevalue");
  } else if (eid === "lhs-button6") {
    tb = document.getElementById("fd6-facevalue");
  } else if (eid === "lhs-button7") {
    tb = document.getElementById("fd7-facevalue");
  } else if (eid === "lhs-button8") {
    tb = document.getElementById("fd8-facevalue");
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
  // var audio = new Audio('https://media.geeksforgeeks.org/wp-content/uploads/20190531135120/beep.mp3');
  if (eid === "lhs-button1") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = true;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
    document.getElementById("lhs-button8").disabled = false;
  } else if (eid === "lhs-button2") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = true;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
    document.getElementById("lhs-button8").disabled = false;
  } else if (eid === "lhs-button3") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = true;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
    document.getElementById("lhs-button8").disabled = false;
  } else if (eid === "lhs-button4") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = true;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
    document.getElementById("lhs-button8").disabled = false;
  // } else if (eid === "lhs-button5") {
  //   // audio.play();
  //   document.getElementById("lhs-button1").disabled = false;
  //   document.getElementById("lhs-button2").disabled = false;
  //   document.getElementById("lhs-button3").disabled = false;
  //   document.getElementById("lhs-button4").disabled = false;
  //   document.getElementById("lhs-button5").disabled = true;
  //   document.getElementById("lhs-button6").disabled = false;
  //   document.getElementById("lhs-button7").disabled = false;
  //   document.getElementById("lhs-button8").disabled = false;
  } else if (eid === "lhs-button6") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = true;
    document.getElementById("lhs-button7").disabled = false;
    document.getElementById("lhs-button8").disabled = false;
  } else if (eid === "lhs-button7") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = true;
    document.getElementById("lhs-button8").disabled = false;
  } else if (eid === "lhs-button8") {
    // audio.play();
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = false;
    document.getElementById("lhs-button4").disabled = false;
    document.getElementById("lhs-button5").disabled = true;
    document.getElementById("lhs-button6").disabled = false;
    document.getElementById("lhs-button7").disabled = false;
    document.getElementById("lhs-button8").disabled = true;
  }
}
