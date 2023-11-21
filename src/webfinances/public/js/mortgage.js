
document.addEventListener('DOMContentLoaded', (event) => {
//  disableElements('{{ .CurrentButton }}');
  disableElements(currentButton);
//  setFocus('{{ .CurrentButton }}');
  setFocus(currentButton);
});

function setFocus(eid) {
  let tb;
var audio = new Audio('https://media.geeksforgeeks.org/wp-content/uploads/20190531135120/beep.mp3');

  if (eid === "lhs-button1") {
    tb = document.getElementById("fd1-n");

 audio.play();

  } else if (eid === "lhs-button2") {
    tb = document.getElementById("fd2-n");
 audio.play();

  } else if (eid === "lhs-button3") {
    tb = document.getElementById("fd3-mrate");
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
  } else if (eid === "lhs-button2") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = true;
    document.getElementById("lhs-button3").disabled = false;
  } else if (eid === "lhs-button3") {
    document.getElementById("lhs-button1").disabled = false;
    document.getElementById("lhs-button2").disabled = false;
    document.getElementById("lhs-button3").disabled = true;
  }
}
