import React from "react";
import { GroupMenu } from "../components/group/GroupMenu";
import { Navigation } from "../components/common/Navigation";

export const Home = ({ GlobalState }) => {
  return (
    <div>
      <Navigation GlobalState={GlobalState} />
      <div className="page-content">Content</div>
      <div className="group-content">
        <GroupMenu GlobalState={GlobalState} />
      </div>
    </div>
  );
};
