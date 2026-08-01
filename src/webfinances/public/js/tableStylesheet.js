/*xxx
Multiple document.addEventListener("DOMContentLoaded", ...) listeners in a single file will run sequentially in the exact order they are
registered (from top to bottom).
*/
document.addEventListener("DOMContentLoaded", function() {
  console.log("\n\nEntering document.addEventListener #1...");
//  const tableBody = document.querySelector("#custom-table tbody");
  const tableBody = document.querySelector("tbody");
  const STORAGE_KEY = "lastSelectedAccount";
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
    const savedId = localStorage.getItem(STORAGE_KEY);
    const allRows = tableBody.querySelectorAll(".clickable-row");
    let rowToSelect = null;
    //Loop through the table rows manually to find a string match.
    if (savedId && allRows.length > 0) {
      for (let i = 0; i < allRows.length; i++) {
        if (allRows[i].getAttribute("data-id") === savedId) {
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
  console.log("Exiting document.addEventListener #1...");
});

//==============================================================
//SECTION 3: Unified Button Click & Secure Form POST Submitter
//==============================================================
document.addEventListener("click", function(event) {
  console.log("\n\nEntering document.addEventListener #2...");
  //Log exactly what element the mouse pointer touched.
  console.log("Mouse pointer physically hit this HTML tag:", event.target);
  //Target check using multiple lookup strategies.
  //Safe check using combined selector logic.
  const clickedButton = event.target.id === "btdelete" ? event.target : event.target.closest("#btdelete");
  if (!clickedButton) return;
  //Prevent premature default browser actions.
  event.preventDefault();
// console.log("Success! Save button detected.");
  console.log("Success! Save button (#btdelete) detected via delegation.");
  // const dynamicTableBody = document.querySelector("#custom-table tbody");
  const dynamicTableBody = document.querySelector("tbody");
  if (!dynamicTableBody) {
    console.error("Table body target #custom-table tbody was not found!");
    return;
  }
  const activeRow = dynamicTableBody.querySelector(".clickable-row.active-row");
  if (!activeRow) {
    alert("Please select a row first.");
    return;
  }
  const id = activeRow.getAttribute("data-id");
  const userConfirmed = confirm(`Are you sure you want to continue with id: "${id}"?`);
  if (userConfirmed) {
//    console.log("User chose CONTINUE. Processing ID:", id);
    console.log("User chose CONTINUE. Packing payload for ID:", id);
    //Locate the form and hidden input fields.
    const hiddenInput = document.getElementById("table_row_id");
    const rowIdForm = document.getElementById("table_id_form");
    if (hiddenInput && rowIdForm) {
      //Load the current active table row ID directly into the hidden element value.
      hiddenInput.value = id;
      console.log("Submitting secure POST request payload containing CSRF token...");
      //Dispatch the form standard POST request to the Go backend endpoint.
      rowIdForm.submit();
    } else {
      console.error("Critical Failure: Form layout elements could not be referenced.");
    }
                // --- YOUR GO ENDPOINT ACTION GOES HERE ---
                // Example A: Redirect to a Go handler URL
                // window.location.href = "/account/process?id=" + encodeURIComponent(accountId);

                // Example B: Submit it to a hidden form field
                // document.getElementById("hidden-account-input").value = accountId;
                // document.getElementById("my-form").submit();
  } else {
    console.log("User chose CANCEL.");
  }
  console.log("Exiting document.addEventListener... #2");
});
