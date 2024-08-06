import React from "react";
import { Link } from "react-router-dom";
import { SubmitForm, getSessionId } from "../../Extra";

export const Navigation = ({ GlobalState }) => {
  const { setSessionUUID } = GlobalState;
  return (
    <nav className="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
      <div className="container-fluid">
        <Link to="/" className="navbar-brand">
          Home
        </Link>
        <form className="d-flex" role="search">
          <input
            className="form-control me-2"
            type="search"
            placeholder="Search"
            aria-label="Search"
          />
          <button className="btn btn-outline-success" type="submit">
            Search
          </button>
        </form>

        <div>
          <ul className="navbar-nav me-auto mb-2 mb-md-0">
            <li className="nav-item">
              <button className="nav-link" href="#">
                Profile
              </button>
            </li>
            <li className="nav-item">
              <button
                className="nav-link"
                onClick={async (e) => {
                  const url = "/logout";
                  const formData = new FormData();
                  formData.append("session", GlobalState.SessionUUID);
                  await SubmitForm(e, url, formData, GlobalState);
                  setSessionUUID(getSessionId());
                }}
              >
                Log out
              </button>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
};
