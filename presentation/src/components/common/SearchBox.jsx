import React, { useState, useEffect } from "react";
import { SubmitForm } from "../../Extra";
import { InputBox } from "../common/InputBox";
export const SearchBox = ({ GlobalState, getValue, url }) => {
  url = "/user_search";
  const [searchData, setSearchData] = useState([]);
  const [isFetching, setIsFetching] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const handleChange = () => {
    console.log("Tere");
    const fetchData = async (e) => {
      const formData = new FormData();
      formData.append("session", GlobalState.SessionUUID);
      const data = await SubmitForm(e, url, formData, GlobalState);
      setSearchData(data.info);
      setIsFetching(false);
      getValue(searchValue);
    };
    fetchData();
  };

  return (
    <div className="form-floating">
      <datalist id="suggestions">
        {searchData.map((x) => {
          console.log({ x });
          return <option>{x}</option>;
        })}
      </datalist>
      <input
        type="text"
        className="form-control height-25rem mb-3"
        id="searchBox"
        name="searchBox"
        placeholder="New member"
        required="required"
        value={searchValue}
        autoComplete="on"
        list="suggestions"
        onChange={(e) => {
          setSearchValue(e.target.value);
          getValue(e.target.value);
          handleChange();
        }}
      />
      <label className="fsize-10" htmlFor="searchBox">
        Find user
      </label>

      {/* <InputBox
        type="text"
        id="1"
        label="Find user"
        requred
        defaulText="New member"
        getValue={(val) => getValue(val)}
        onChange={(e) => {
          setSearchValue(e.target.value);
          getValue(e.target.value);
          handleChange();
        }}
        autoComplete="on"
        list="suggestions"
      /> */}
    </div>
  );
};
