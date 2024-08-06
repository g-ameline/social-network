import React, { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { Home } from "./pages/Home";
import { ErrorPrint } from "./components/common/Error";
import { getSessionId } from "./Extra";

function App() {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isError, setIsError] = useState("");
  const [appUser, setAppUser] = useState("");
  const [SessionUUID, setSessionUUID] = useState("");
  const [groups, setGroups] = useState([]);
  const GlobalState = {
    isSubmitting,
    setIsSubmitting,
    isError,
    setIsError,
    appUser,
    setAppUser,
    SessionUUID,
    setSessionUUID,
    groups,
    setGroups,
  };

  useEffect(() => {
    const cookie = getSessionId();
    setSessionUUID(cookie);
  }, []);

  return (
    <div className="App">
      {isError !== "" && <ErrorPrint content={isError} />}
      <BrowserRouter>
        <Routes>
          <Route
            exact
            path="/"
            element={
              SessionUUID !== "" ? (
                <Home GlobalState={GlobalState} />
              ) : (
                <Navigate replace to={"/login"} />
              )
            }
          />
          <Route
            exact
            path="/login"
            element={
              SessionUUID === "" ? (
                <Login GlobalState={GlobalState} />
              ) : (
                <Navigate replace to={"/"} />
              )
            }
          />
          <Route
            exact
            path="/register"
            element={
              SessionUUID === "" ? (
                <Register GlobalState={GlobalState} />
              ) : (
                <Navigate replace to={"/"} />
              )
            }
          />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
