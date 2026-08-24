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
      if (allRows[i].getAttribute("data-id") === savedRowId) {
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
    const currentId = clickedRow.getAttribute("data-id");
    localStorage.setItem(STORAGE_KEY, currentId);
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
      const selectedId = activeRow.getAttribute("data-id");
      const userConfirmed = confirm(`Are you sure you want to continue with id: "${selectedId}"?`);
      if (userConfirmed) {
        const hiddenInputId = document.getElementById("table_row_id");
        const hiddenInputStyle = document.getElementById("hidden_bttable");
        if (hiddenInputId && hiddenInputStyle) {
          //Sync runtime data to hidden inputs for backend delivery
          hiddenInputId.value = selectedId;
          hiddenInputStyle.value = saveButton.value;
          //Dispatch secure payload to Go Backend
          mainForm.submit();
        } else {
          console.error("Form target fields are completely missing from the HTML DOM.");
        }
      }
    });
  }
  console.log("Exiting (1) tableStylesheet.js...");
});
