import { useState, type FC } from "react";
import type { RootState } from "../../store/store";
import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";
import { Button, Form } from "react-bootstrap";
import "./new-password.css";
import { useTranslation } from "react-i18next";
import { recoveryPassword } from "../../api/auth-api/auth-api";
import type { RecoveryPasswordData } from "../../../types/types";

const NewPassword: FC = () => {
  const { t } = useTranslation();
  const tokenRecoveryPassword = useSelector(
    (state: RootState) => state.auth.tokenRecoveryPassword.token,
  );
  const email = useSelector((state: RootState) => state.auth.emailRecoveryPassword);
  const [newPassword, setNewPassword] = useState("")
  const [passwordConfirmation, setPasswordConfirmation] = useState("")

  if (!tokenRecoveryPassword) {
    return <Navigate to={"/"} />;
  }

  function handleRecoveryPassword() {
    console.log(123)
    const dataRecoveryPassword: RecoveryPasswordData = {
      token: tokenRecoveryPassword,
      phone: "",
      email: email ? email : "",
      new_password: newPassword,
      confirmation_password: passwordConfirmation
    }
    recoveryPassword(dataRecoveryPassword).then((data) => console.log(data))
  }

  return (
    <div className="new-password-page">
      <div className="new-password-wrapper">
        <div className="new-password-block">
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
            <Button onClick={handleRecoveryPassword}>
              {t("auth.send")}
            </Button>
        </div>
      </div>
    </div>
  );
};

export default NewPassword;
