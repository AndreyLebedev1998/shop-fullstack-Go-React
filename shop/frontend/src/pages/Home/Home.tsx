import { type FC } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../../store/store";
import AlreadyBought from "../../components/AlreadyBought/AlreadyBought";
import Recommendations from "../../components/Recommendations/Recommendations";
import "./home.css";

const Home: FC = () => {
  const token = useSelector((state: RootState) => state.auth.token);
  return (
    <>
      <div className="home-page">
        {token && <div className="already-bought-wrapper">
          <AlreadyBought/>
        </div>}
        {token && <div className="recommendations-wrapper">
          <Recommendations />
        </div>}
      </div>
    </>
  );
};

export default Home;
