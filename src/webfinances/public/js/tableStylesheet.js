/*
Multiple document.addEventListener("DOMContentLoaded", ...) listeners in a single file will run sequentially in the exact order they are
registered (from top to bottom).
*/
document.addEventListener("DOMContentLoaded", function() {
  console.log("\n\nEntering tableStylesheet.js #1...");
  const tableBody = document.querySelector("tbody");
  const STORAGE_KEY = "lastSelectedRow";
  //==========================================================
  //SECTION 1: Handle Table Row Click Interactions & Memory
  //==========================================================
  if (tableBody) {
    tableBody.addEventListener("click", function(event) {
      //Find the closest table row that was clicked.
      const row = event.target.closest(".clickable-row");
      //If the click wasn't inside a valid table row, do nothing.
      if (!row) return;
      //Remove selection class from all rows.
      const allRows = tableBody.querySelectorAll(".clickable-row");
      allRows.forEach(r => r.classList.remove("active-row"));
      //Add selection class to the clicked row.
      row.classList.add("active-row");
      //Extract and save the id.
      const id = row.getAttribute("data-id");
      //Console log it so you can see it working in F12 Inspect tools.
      console.log("Selected id is:", id);
      //SAVE to browser memory.
      localStorage.setItem(STORAGE_KEY, id);
    });
    //==========================================================
    //SECTION 2: Restore Memory Highlight or Default to Row #1
    //==========================================================
    /*
    You do not need to manually clear the memory because it is automatically overridden every single time a user clicks a different row.
    localStorage works like a single variable that holds one value at a time. When you call localStorage.setItem(STORAGE_KEY, variableName),
     it instantly deletes the old value and writes the new one over it.
    */
    //Restore Memory or Fallback to First Row. Try to read the saved id from localStorage.
    const savedRow = localStorage.getItem(STORAGE_KEY);
    const allRows = tableBody.querySelectorAll(".clickable-row");
    let rowToSelect = null;
    //Loop through the table rows manually to find a string match.
    if (savedRow && allRows.length > 0) {
      for (let i = 0; i < allRows.length; i++) {
        if (allRows[i].getAttribute("data-id") === savedRow) {
          rowToSelect = allRows[i];
          break;  //Match found! Stop searching.
        }
      }
    }
    //If nothing was saved before or the id no longer exists, default to the first row.
    if (!rowToSelect && allRows.length > 0) {
      rowToSelect = allRows[0];
    }
    //Physically apply the active class if a valid row was targetable.
    if (rowToSelect) {
      rowToSelect.classList.add("active-row");
      console.log("Active row applied successfully:", rowToSelect.getAttribute("data-id"));
    }
  }
  console.log("Exiting tableStylesheet.js #1...");
});

//==============================================================
//SECTION 3: Unified Button Click & Secure Form POST Submitter
//==============================================================
/*
When you submit a form using the JavaScript's form.submit() method, the browser only sends the data from standard input elements (like <input>,
<select>, or <textarea>). It completely ignores the name and value attributes of <button> elements.

When a human physically clicks a standard submit button (type="submit"), the browser includes that specific button's name/value because it was
the trigger. However, when you change a button to type="button" and force a submission via JavaScript (form.submit()), the browser treats the
form as an isolated object and only serializes actual form fields (<input>). Using hidden inputs completely bypasses this limitation.
*/
document.addEventListener("click", function(event) {
  console.log("\n\nEntering tableStylesheet.js #2...");
  //Log exactly what element the mouse pointer touched.
  console.log("Mouse pointer physically hit this HTML tag:", event.target);
  //Target check using multiple lookup strategies.
  //Safe check using combined selector logic.
  const clickedButton = event.target.id === "bttable" ? event.target : event.target.closest("#bttable");
  if (!clickedButton) {
    console.error("Cannot find the clicked button.");
    console.log("Exiting tableStylesheet.js... #2");
    return;
  }
  //Prevent premature default browser actions.
  event.preventDefault();
  console.log("Success! Save button detected via delegation:", clickedButton);
  const dynamicTableBody = document.querySelector("tbody");
  if (!dynamicTableBody) {
    console.error("Table body target #custom-table tbody was not found!");
    console.log("Exiting tableStylesheet.js... #2");
    return;
  }
  const activeRow = dynamicTableBody.querySelector(".clickable-row.active-row");
  if (!activeRow) {
    alert("Please select a row first.");
    console.log("Exiting tableStylesheet.js... #2");
    return;
  }
  const id = activeRow.getAttribute("data-id");
  const userConfirmed = confirm(`Are you sure you want to continue with id: "${id}"?`);
  if (userConfirmed) {
    console.log("CONTINUE. Processing ID:", id);
    //Locate the form and hidden input fields.
    const hiddenInput = document.getElementById("table_row_id");
    const tableIdForm = document.getElementById("table_id_form");
    if (hiddenInput && tableIdForm) {
      //Load the current active table row ID directly into the hidden element value.
      hiddenInput.value = id;
      const hiddenInputButtonId = document.getElementById("hidden_bttable");
      //Grab the value directly from the clicked button and put it in the hidden input
      if (hiddenInputButtonId) {
        hiddenInputButtonId.value = clickedButton.value; // This grabs "rhs-*"
        console.log("Loaded button value into payload:", clickedButton.value);
      }
      console.log("Submitting secure POST request payload...");
      //Dispatch the form standard POST request to the Go backend endpoint.
      tableIdForm.submit();
    } else {
      console.error("Critical Failure: Form layout elements could not be referenced.");
    }
  } else {
    console.log("CANCEL.");
  }
  console.log("Exiting tableStylesheet.js... #2");
});
