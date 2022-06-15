import toast from "react-hot-toast";

export const API_URL = process.env.REACT_APP_API_URL;

function handleNon2XXCode(response) {
  if (response.ok) {
    return response;
  }
  return response
    .text()
    .catch((x) => Promise.reject("failed to load"))
    .then((text) =>
      Promise.reject(`${text}: response code: ${response.status}`)
    );
}

function toastErrors(defaultValue) {
  return function(err) {
    toast.error(err.toString(), {
      duration: 10000,
    });
    return defaultValue;
  };
}

export function apiGetCategories() {
  return fetch(`${API_URL}/api/category`)
    .then(handleNon2XXCode)
    .then((res) => res.json())
    .catch(toastErrors([]));
}

export function apiGetCategory(id) {
  return fetch(`${API_URL}/api/category/${id}`)
    .then(handleNon2XXCode)
    .then((res) => res.json())
    .catch(toastErrors());
}

export function apiImageUrl(fileName) {
  return `${API_URL}/images/${fileName}`;
}
