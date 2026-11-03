// Filter forms: <form data-filter-form hx-trigger="filter-change" …> refetches only
// when its named values change.
//
// hx-trigger="change" is not enough on a form that mixes controls. A combobox's
// visible text input has no name, yet the browser fires a native change on it when
// it loses focus, for example on the click that opens a sheet in the results, and the
// refetch replaces the results under that click. htmx's `changed` modifier doesn't
// help: it compares the triggering element's own value, and datepicker.js fires
// change on the form itself.
//
// So: on any change inside the form (or on the form), compare its form data with what
// it last sent and fire "filter-change" only when it differs. Nothing is stopped, so
// components that listen for change keep working.
//
// "What it last sent" starts as the values the page loaded with, recorded the first
// time focus enters the form, which comes before any change a user can make.
// Without it the first change always fired, so the blur case above still refetched
// once per page load.
const snapshot = (form) => new URLSearchParams(new FormData(form)).toString();

document.addEventListener("focusin", (e) => {
  const form = e.target instanceof Element && e.target.closest("form[data-filter-form]");
  if (form && form._filterValues === undefined) form._filterValues = snapshot(form);
});

document.addEventListener("change", (e) => {
  const form = e.target instanceof Element && e.target.closest("form[data-filter-form]");
  if (!form) return;
  const values = snapshot(form);
  if (values === form._filterValues) return;
  form._filterValues = values;
  form.dispatchEvent(new CustomEvent("filter-change"));
});
