async function loadNextPage() {
  const sentinel = document.getElementById("pagination-sentinel");
  if (!sentinel) return;

  const targetCell = sentinel.querySelector("td");
  if (!targetCell) return;

  const cursor = targetCell.dataset.nextCursor;
  const hasNext = targetCell.dataset.hasNext === "true";
  if (!hasNext || !cursor) return;

  const csrfInput = document.querySelector('input[name="csrf_token"]');
  const csrfToken = csrfInput ? csrfInput.value : '';

  const requestBody = new URLSearchParams();
  requestBody.append("csrf_token", csrfToken);
  requestBody.append("bttable", "load-more");
  requestBody.append("cursor", cursor);

  // Target the exact URL path you are currently standing on
  const targetRoute = window.location.pathname;

  // ==========================================================================
  // CACHE BUSTER: Append a unique timestamp to force the browser to skip cache
  // ==========================================================================
  const cacheBusterUrl = `${targetRoute}?_cb=${new Date().getTime()}`;

  //Dispatch the POST action silently in the background using the unique URL
  const response = await fetch(cacheBusterUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    body: requestBody.toString()
  });

  if (!response.ok) {
    console.error("Asynchronous row append failed.");
    return;
  }

  const htmlFragment = await response.text();
  sentinel.remove();

  const tbody = document.getElementById("tbody-users");
  if (tbody) {
    tbody.insertAdjacentHTML("beforeend", htmlFragment);
  }

  const newSentinelCell = document.querySelector("#pagination-sentinel td");
  const loadMoreBtn = document.getElementById("load-more-btn");

  if (loadMoreBtn && newSentinelCell) {
    if (newSentinelCell.dataset.hasNext === "true") {
      loadMoreBtn.removeAttribute("disabled"); // Ensure button isn't locked if data continues
      loadMoreBtn.classList.remove("hidden");
    } else {
      loadMoreBtn.setAttribute("disabled", "true"); // Explicitly disable button if no records remain
      loadMoreBtn.classList.add("hidden");
    }
  }

  // =========================================================================
  // STATE PROTECTION: Reset form fields and disable Delete button on append
  // =========================================================================
  const deleteBtn = document.getElementById('bttable');
  const hiddenBtTable = document.getElementById("hidden_bttable");
  const hiddenId = document.getElementById("hidden_id");

  if (deleteBtn) deleteBtn.setAttribute("disabled", "true");
  if (hiddenBtTable) hiddenBtTable.value = "";
  if (hiddenId) hiddenId.value = "";

  // Clear any leftover visual highlights on older rows to avoid confusion
  document.querySelectorAll('.clickable-row').forEach(row => {
    row.classList.remove('active-row');
  });
}
