import React, { useState } from "react";
import { CloseFreezeButton } from "../common/CloseFreezeButton";
import { ImageSelect } from "../common/ImageSelect";
import { InputBox } from "../common/InputBox";
import { SubmitButton } from "../common/SubmitButton";
import { SubmitForm } from "../../Extra";

export const GroupAdd = ({ title, GlobalState }) => {
  const url = "/group_create";
  const css = "add-right-pos";
  const [isClick, SetIsClicked] = useState(false);
  const [, setLabel] = useState();
  const [, setDesc] = useState();
  const [, setImage] = useState("");
  const handleSubmit = async (e) => {
    await SubmitForm(e, url, [], GlobalState);
    SetIsClicked(false);
  };
  return (
    <>
      <button
        title={title}
        className={
          "bi bi-plus-square btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed " +
          css
        }
        onClick={(e) => {
          SetIsClicked(true);
        }}
      ></button>
      {isClick === true && (
        <div className="freezePanel">
          <CloseFreezeButton getValue={(val) => SetIsClicked(val)} />
          <main className="form-parent">
            <form onSubmit={handleSubmit}>
              <ImageSelect getValue={(img) => setImage(img)} />
              <InputBox
                type="text"
                id="name"
                label="Group name"
                defaulText="Group name"
                required
                getValue={(val) => setLabel(val)}
              />
              <InputBox
                type="text"
                id="dexcroption"
                label="Group description"
                defaulText="Group description"
                required
                getValue={(val) => setDesc(val)}
              />
              <SubmitButton text="Create new Group" />
            </form>
          </main>
        </div>
      )}
    </>
  );
};
