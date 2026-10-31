document.addEventListener("DOMContentLoaded", function () {
  console.log("\nEntering tableStylesheet.js...");
  const tableBody = document.getElementById("tbody");
  const saveButton = document.getElementById("bttable");
  const mainForm = document.getElementById("table_id_form");
  const STORAGE_KEY = "lastSelectedRow";
  if (!tableBody) {
    console.log("Exiting (2) tableStylesheet.js...");
    return;
  }
  const allRows = tableBody.querySelectorAll(".clickable-row");
  // =========================================================================
  // MEMORY MANAGEMENT (Restore Highlight on Load or Default to Row 1)
  // =========================================================================
  const savedRowId = localStorage.getItem(STORAGE_KEY);
  let rowToSelect = null;
  if (savedRowId && allRows.length > 0) {
    //Find the row matching our saved localStorage item
    for (let i = 0; i < allRows.length; i++) {
      if (allRows[i].getAttribute("data-row") === savedRowId) {
        rowToSelect = allRows[i];
        break;
      }
    }
  }
  //Fallback default: Select row #1 if nothing was stored or found
  if (!rowToSelect && allRows.length > 0) {
    rowToSelect = allRows[0];
  }
  //
  if (rowToSelect) {
    rowToSelect.classList.add("active-row");
  }
  // =========================================================================
  // ROW CLICK LISTENERS (Switch Highlights & Store Selection)
  // =========================================================================
  tableBody.addEventListener("click", function (event) {
    const clickedRow = event.target.closest(".clickable-row");
    if (!clickedRow) {
      return;
    }
    //Clear old active highlight selections
    allRows.forEach(row => row.classList.remove("active-row"));
    //Set new active highlight row selection
    clickedRow.classList.add("active-row");
    //Save selection to persistent memory
    const currentRow = clickedRow.getAttribute("data-row");
    localStorage.setItem(STORAGE_KEY, currentRow);
  });
  // =========================================================================
  // SECURE DATA SUBMISSION (Delete User Click Workflow)
  // =========================================================================
  if (saveButton && mainForm) {
    saveButton.addEventListener("click", function (event) {
      event.preventDefault();  //Stop instant submission cycles
      const activeRow = tableBody.querySelector(".clickable-row.active-row");
      if (!activeRow) {
        alert("Please select a row first.");
        return;
      }
      const selectedRow = activeRow.getAttribute("data-row");
      const selectedId = activeRow.getAttribute("data-id");
      //Verify that the data attributes actually exist before showing the confirm dialog
      if (!selectedRow || !selectedId) {
        console.error("The selected row is missing data-row or data-id attributes.");
        return;
      }
      const userConfirmed = confirm(`Are you sure you want to continue with row: "${selectedRow}"?`);
      if (userConfirmed) {
        // const hiddenInputRow = document.getElementById("table_row_id");
        const hiddenInputStyle = document.getElementById("hidden_bttable");
        const hiddenInputId = document.getElementById("hidden_id");
        //Comprehensive safety check to ensure all elements exist in the DOM
        if (/*hiddenInputRow &&*/ hiddenInputStyle && hiddenInputId) {
          //Sync runtime data to hidden inputs for backend delivery.
          // hiddenInputRow.value = selectedRow.trim();
          hiddenInputStyle.value = saveButton.value.trim();
          hiddenInputId.value = selectedId.trim();
          //Guard: Ensure we aren't submitting empty identifier keys to Go
          if (/*hiddenInputRow.value !== "" &&*/ hiddenInputId.value !== "") {
            //Dispatch payload to Go Backend cleanly.
            mainForm.submit();
          } else {
            console.warn("Form submission halted: row data or unique identifier is blank.");
          }
        } else {
          console.error("Form target fields are completely missing from the HTML DOM.");
        }
      }
    });
  }
  console.log("Exiting (1) tableStylesheet.js...");
});
