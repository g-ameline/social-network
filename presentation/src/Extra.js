export const getSessionId = () => {
  let cookie = document.cookie;
  let cookies = cookie.split("; ");
  cookie = cookies.find((row) => row.startsWith("Social-network="));
  return typeof cookie === "undefined"
    ? ""
    : cookie.substring(15, cookie.length);
};

export const SubmitForm = async (e, url, formContent, GlobalState) => {
  if (typeof e !== "undefined") {
    e.preventDefault();
  }
  console.log("Send to server:", url, formContent, e);
  const { setIsSubmitting, setIsError, setSessionUUID } = GlobalState;
  setIsError("");
  setIsSubmitting(true);
  const formData =
    formContent.length === 0 ? new FormData(e.target) : formContent;
  try {
    const response = await fetch(url, {
      method: "POST",
      body: formData,
    });
    const result = await response.json();
    console.log("result:", result);
    if (result.what === "error") {
      setIsSubmitting(false);
      setIsError(result.info);
      return;
    }
    setSessionUUID(getSessionId());
    setIsSubmitting(false);
    return result;
  } catch (error) {
    setIsSubmitting(false);
    setIsError(error);
    return;
  }
};
