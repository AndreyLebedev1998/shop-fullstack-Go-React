import { useState, type FC } from "react";
import type { RootState } from "../../store/store";
import { useSelector } from "react-redux";
import { Navigate, useNavigate } from "react-router-dom";
import { Alert, Button, Form, OverlayTrigger, Tooltip } from "react-bootstrap";
import "./new-password.css";
import { useTranslation } from "react-i18next";
import { recoveryPassword } from "../../api/auth-api/auth-api";
import type { RecoveryPasswordData } from "../../../types/types";

const NewPassword: FC = () => {
  const navigate = useNavigate()
  const { t } = useTranslation();
  const tokenRecoveryPassword = useSelector(
    (state: RootState) => state.auth.tokenRecoveryPassword.token,
  );
  const email = useSelector(
    (state: RootState) => state.auth.emailRecoveryPassword,
  );
  const [newPassword, setNewPassword] = useState("");
  const [passwordConfirmation, setPasswordConfirmation] = useState("");
  const [show, setShow] = useState(false);
  const [notMatchPassword, setNotMatchPassword] = useState(false);
  const [passwordLessSymbols, setPasswordLessSymbols] = useState(false);
  const [error, setError] = useState(false)
  const [successRecoveryPassword, setSuccessRecoveryPassword] = useState(false)

  if (!tokenRecoveryPassword) {
    return <Navigate to={"/"} />;
  }

  function handleRecoveryPassword() {
    if (show && (notMatchPassword || passwordLessSymbols)) return;
    if (newPassword !== passwordConfirmation) {
      setShow(true);
      setNotMatchPassword(true);
      return;
    }

    if (newPassword.length < 5) {
      setShow(true);
      setPasswordLessSymbols(true);
      return;
    }

    const dataRecoveryPassword: RecoveryPasswordData = {
      token: tokenRecoveryPassword,
      phone: "",
      email: email ? email : "",
      new_password: newPassword,
      confirmation_password: passwordConfirmation,
    };
    recoveryPassword(dataRecoveryPassword, setError).then((data) => {
      if (data?.message) {
        setSuccessRecoveryPassword(true)
        setTimeout(() => {
          navigate("/authorization")
        }, 1500)
      }
    });
  }

  return (
    <div className="new-password-page">
      <div className="new-password-wrapper">
        <div className="new-password-block">
          {error && (
            <Alert
              className="error-alert"
              onClose={() => setError(false)}
              dismissible
              variant="danger"
            >
              {t("auth.error")}
            </Alert>
          )}
          {successRecoveryPassword && (
            <Alert
              className="success-recovery-password"
              onClose={() => setSuccessRecoveryPassword(false)}
              dismissible
              variant="success"
            >
              {t("auth.success_recovery_password")}
            </Alert>
          )}
          <Form.Label htmlFor="inputPassword5">{t("auth.password")}</Form.Label>
          <Form.Control
            type="password"
            id="inputPassword5"
            aria-describedby="passwordHelpBlock"
            onChange={(e) => setNewPassword(e.target.value)}
          />
          <Form.Text id="passwordHelpBlock" muted>
            {t("auth.new_password")}
          </Form.Text>
        </div>
        <div className="new-confirmation-password-block">
          <Form.Label htmlFor="inputPassword5">
            {t("auth.password_confirmation")}
          </Form.Label>
          <Form.Control
            type="password"
            id="inputPassword5"
            onChange={(e) => setPasswordConfirmation(e.target.value)}
          />
        </div>
        <div className="new-password-action">
          <OverlayTrigger
            placement="top"
            show={show && (notMatchPassword || passwordLessSymbols)}
            overlay={
              notMatchPassword ? (
                <Tooltip>{t("auth.not_match_password")}</Tooltip>
              ) : (
                <Tooltip>{t("auth.new_password")}</Tooltip>
              )
            }
          >
            <Button
              onClick={() => {
                handleRecoveryPassword();
                setTimeout(() => {
                  setShow(false)
                  setNotMatchPassword(false)
                  setPasswordLessSymbols(false)
                }, 3000);
              }}
            >
              {t("auth.send")}
            </Button>
          </OverlayTrigger>
        </div>
      </div>
    </div>
  );
};

export default NewPassword;
