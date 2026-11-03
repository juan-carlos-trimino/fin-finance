
document.addEventListener("DOMContentLoaded", function () {
  console.log("Entering tableStylesheet.js...");
  const tableBody = document.getElementById("table-body");
  const actionButton = document.getElementById("action-button");
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
      if (allRows[i].getAttribute("data-username") === savedRowId) {
        rowToSelect = allRows[i];
        break;
      }
    }
  }
  // Only add the visual selection class if a row was actually found in memory
  if (rowToSelect) {
    rowToSelect.classList.add("active-row");
    //Enable the delete button since an active row has been recovered
    if (actionButton) {
      actionButton.removeAttribute("disabled");
    }
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
    const currentRow = clickedRow.getAttribute("data-username");
    localStorage.setItem(STORAGE_KEY, currentRow);
  });
  // =========================================================================
  // SECURE DATA SUBMISSION (Delete User Click Workflow)
  // =========================================================================
  if (actionButton && mainForm) {
    actionButton.addEventListener("click", function (event) {
      event.preventDefault();  //Stop instant submission cycles
      const activeRow = tableBody.querySelector(".clickable-row.active-row");
      if (!activeRow) {
        alert("Please select a row first.");
        return;
      }
      const selectedRow = activeRow.getAttribute("data-username");
      const selectedId = activeRow.getAttribute("data-id");
      //Verify that the data attributes actually exist before showing the confirm dialog
      if (!selectedRow || !selectedId) {
        console.error("The selected row is missing data-row or data-id attributes.");
        return;
      }
      const userConfirmed = confirm(`Are you sure you want to continue with row: "${selectedRow}"?`);
      if (userConfirmed) {
        const hiddenInputStyle = document.getElementById("hidden_action-button");
        const hiddenInputId = document.getElementById("hidden_id");
        //Comprehensive safety check to ensure all elements exist in the DOM
        if (hiddenInputStyle && hiddenInputId) {
          //Sync runtime data to hidden inputs for backend delivery.
          hiddenInputStyle.value = "delete";
          hiddenInputId.value = selectedId.trim();
          //Guard: Ensure we aren't submitting empty identifier keys to Go
          if (hiddenInputId.value !== "") {
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

/***
Selecting the asynchronous background append via fetch is the ideal choice for a smooth, modern UI. It prevents the browser from
refreshing the full page, keeps the existing table rows on the screen, and appends the new records seamlessly into your container.

Since we are replacing the form.submit() approach with a background network request, you must make sure that Go returns only the
raw <tr> table rows back to the browser during a POST pagination request. If Go tries to return the full page layout, your table
cells will break and text will distort.
***/
document.addEventListener('DOMContentLoaded', () => {
  console.log("Entering document.addEventListener1...");
  const tbody = document.getElementById('table-body');
  if (!tbody) return;

  function handleRowSelection(rowElement) {
    // 1. Extract values from your updated 'data-username' and 'data-id' attributes
    const userName = rowElement.getAttribute('data-username');
    const id = rowElement.getAttribute('data-id');
    console.log(`Row selected securely: Username=${userName}, ID=${id}`);

    // 2. Locate your existing hidden form input elements
    const hiddenBtTable = document.getElementById("hidden_action-button");
    const hiddenId = document.getElementById("hidden_id");
    if (hiddenBtTable && hiddenId) {
      // Drop the values directly into the hidden form buckets
      hiddenBtTable.value = "delete";
      hiddenId.value = id;  // This gives the "Delete User" button the target primary key ID.
      console.log(`Hidden form buckets updated: action-button=${hiddenBtTable.value}, id=${hiddenId.value}`);
    }

    // 3. VISUAL HIGHLIGHT: Matches your external stylesheet's row state (.active-row)
    document.querySelectorAll('.clickable-row').forEach(row => {
      row.classList.remove('active-row');
    });
    rowElement.classList.add('active-row');

    // 4. UNLOCK DELETE BUTTON: Enable the button now that an active target ID is confirmed
    const deleteBtn = document.getElementById('action-button');
    if (deleteBtn) {
      deleteBtn.removeAttribute('disabled');
    }
  }

  // 1. Capture Mouse Click Events via Delegation
  tbody.addEventListener('click', (event) => {
    const row = event.target.closest('.clickable-row');
    if (row && tbody.contains(row)) {
      handleRowSelection(row);
    }
  });

  // 2. Capture Accessible Keyboard Events via Delegation
  tbody.addEventListener('keydown', (event) => {
    if (event.key === 'Enter' || event.key === ' ') {
      const row = event.target.closest('.clickable-row');
      if (row && tbody.contains(row)) {
        event.preventDefault();  // Stop page scrolling default on spacebar
        handleRowSelection(row);
      }
    }
  });
});

/***
This is a standard defensive programming rule. If a user visits a different page in the dashboard that includes this global script
file but does not display the users data grid (meaning the "load-more-btn" element is completely absent from the current DOM tree),
this conditional check gracefully skips the binding block rather than throwing a fatal 'Uncaught TypeError: Cannot read properties
of null' crash error.
***/
document.addEventListener('DOMContentLoaded', () => {
  // Bind the load-more click handler securely without breaking CSP
  const loadMoreBtn = document.getElementById('load-more-btn');
  if (loadMoreBtn) {
    loadMoreBtn.addEventListener('click', (event) => {
      event.preventDefault();  // Stop the enclosing <form> from triggering a full page reload!
      loadNextPage();  // Execute a smooth, asynchronous row injection script instead
    });
  }
});
