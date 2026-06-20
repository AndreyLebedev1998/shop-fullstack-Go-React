import { useState, type FC } from "react";
import { Alert, Button, Form, Tab, Tabs } from "react-bootstrap";
import InputGroup from "react-bootstrap/InputGroup";
import { useTranslation } from "react-i18next";
import "./authorization.css";
import { authorization, getMe } from "../../api/auth-api/auth-api";
import type { RootState } from "../../store/store";
import Registration from "../Registration/Registration";
import { useDispatch, useSelector } from "react-redux";
import { setAuth } from "../../store/authSlice";
import { Link, Navigate } from "react-router-dom";

const Authorization: FC = () => {
  const dispatch = useDispatch();
  const token = useSelector((state: RootState) => state.auth.token);
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState(false);
  const [errorServer, setErrorServer] = useState(false);
  const [show, setShow] = useState(false);
  const [key, setKey] = useState("auth");
  const [showSuccess, setShowSuccess] = useState(false);
  const { t } = useTranslation();

  if (token) {
    return <Navigate to={"/"} />;
  }

  function submit() {
    if (!email.length || !password.length) return;
    authorization(email, password, setErrorServer, setShow).then((data) => {
      if (!data?.token) {
        setError(true);
        setShow(true);
        return;
      }
      getMe(data.token).then((dataAuth) => {
        if (dataAuth?.token) {
          localStorage.setItem("token", dataAuth?.token);
          dispatch(setAuth(dataAuth));
        }
      });
    });
  }

  return (
    <>
      {showSuccess && (
        <Alert
          className="reg_success_alert"
          onClose={() => setShowSuccess(true)}
          dismissible
          variant="success"
        >
          {t("auth.reg_success")}
        </Alert>
      )}
      <div className="auth-page">
        <div className="auth-page__auth-wrapper">
          {error && show && key !== "reg" && (
            <Alert onClose={() => setShow(false)} dismissible variant="danger">
              {t("auth.error_auth")}
            </Alert>
          )}
          {errorServer && show && key !== "reg" && (
            <Alert onClose={() => setShow(false)} dismissible variant="danger">
              {t("auth.server_error")}
            </Alert>
          )}
          <Tabs
            activeKey={key}
            defaultActiveKey="auth"
            onSelect={(k) => {
              if (k) {
                setKey(k);
              }
            }}
          >
            <Tab eventKey="auth" title={t("auth.auth")}>
              <div className="auth-wrap">
                <InputGroup className="mb-3">
                  <InputGroup.Text id="basic-addon1">@</InputGroup.Text>
                  <Form.Control
                    name="email"
                    autoComplete="email"
                    placeholder={t("auth.email")}
                    onChange={(e) => setEmail(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") {
                        submit();
                      }
                    }}
                  />
                </InputGroup>
              </div>
              <InputGroup className="mb-3">
                <Form.Control
                  placeholder={t("auth.password")}
                  type="password"
                  minLength={5}
                  onChange={(e) => setPassword(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      submit();
                    }
                  }}
                />
              </InputGroup>
              <Link to="/recovery-password">
                  {t("auth.forgot_your_password")}
              </Link>
              <div className="auth-page__action">
                <Button onClick={submit} variant="primary">
                  Войти
                </Button>
              </div>
            </Tab>
            <Tab eventKey="reg" title={t("auth.reg")}>
              <Registration setKey={setKey} setShowSuccess={setShowSuccess} />
            </Tab>
          </Tabs>
        </div>
      </div>
    </>
  );
};

export default Authorization;
