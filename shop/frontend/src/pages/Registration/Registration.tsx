import { useState, type FC } from "react";
import { Alert, Button, Form, InputGroup } from "react-bootstrap";
import { IMaskInput } from "react-imask";
import { useTranslation } from "react-i18next";
import "./regestration.css";
import type { RegRequestData } from "../../../types/types";
import { registration } from "../../api/auth-api/auth-api";

type TypeRegistration = {
    setKey: (key: string) => void
    setShowSuccess: (show: boolean) => void
}

const Registration: FC<TypeRegistration> = ({setKey, setShowSuccess}) => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [name, setName] = useState<string>("");
  const [lastname, setLastName] = useState<string>("");
  const [phone, setPhone] = useState<string>("");
  const [showReqFields, setShowReqFields] = useState(false);
  const [showInncorectFields, setShowInncorectFields] = useState(false);
  const [showSmallPassword, setShowSmallPassword] = useState(false)
  const [error, setError] = useState(false)
  const [errorConflict, setErrorConflict] = useState(false)

  const submit = () => {
    if (!email.length || !password.length || !name.length) {
      setShowReqFields(true);
      return;
    } else {
        setShowReqFields(false);
    }

    if (
      !email.includes("@") ||
      phone.length && phone.length != 18 ||
      name.length < 3
    ) {
    setShowInncorectFields(true)
      return;
    } else {
        setShowInncorectFields(false)
    }

    if (password.length < 5) {
        setShowSmallPassword(true)
        return;
    }else {
        setShowSmallPassword(false)
    }

    const regData: RegRequestData = {
        email,
        name,
        lastname,
        phone,
        password
    }

    registration(regData).then((data) => {
    if (!data) return

    console.log(data)
    
    if ('error' in data) {
        setErrorConflict(true)
        return
    }

    if (data.user_id) {
        setShowSuccess(true)
        setEmail("")
        setPassword("")
        setLastName("")
        setPhone("")
        setName("")
        setKey("auth")
    } else {
        setError(true)
    }
})
  };
  const { t } = useTranslation();
  return (
    <div className="registration-page__registration-wrapper">
      {showReqFields && (
        <Alert onClose={() => setShowReqFields(false)} dismissible variant="warning">
          {t("auth.required_fields")}
        </Alert>
      )}
      {showInncorectFields && (
        <Alert onClose={() => setShowInncorectFields(false)} dismissible variant="warning">
          {t("auth.reg_incorrect_fileds")}
        </Alert>
      )}
      {showSmallPassword && (
        <Alert onClose={() => setShowSmallPassword(false)} dismissible variant="warning">
          {t("auth.small_password")}
        </Alert>
      )}
    {error && (
        <Alert onClose={() => setError(false)} dismissible variant="danger">
          {t("auth.error_reg")}
        </Alert>
      )}
      {errorConflict && (
        <Alert onClose={() => setErrorConflict(false)} dismissible variant="danger">
          {t("auth.email_or_password_is_busy")}
        </Alert>
      )}
      <InputGroup className="mb-3">
        <span className="required__field required__field-reg">*</span>
        <InputGroup.Text id="basic-addon1">@</InputGroup.Text>
        <Form.Control
          required
          name="email"
          autoComplete="email"
          placeholder={t("auth.email")}
          onChange={(e) => setEmail(e.target.value)}
          value={email}
        />
      </InputGroup>
      <InputGroup className="mb-3">
        <span className="required__field required__field-reg">*</span>
        <Form.Control
          required
          name="name"
          autoComplete="name"
          placeholder={t("auth.name")}
          onChange={(e) => setName(e.target.value)}
          value={name}
        />
      </InputGroup>
      <InputGroup className="mb-3 not-required__field-reg">
        <Form.Control
          required
          name="lastname"
          autoComplete="lastname"
          placeholder={t("auth.lastname")}
          onChange={(e) => setLastName(e.target.value)}
          value={lastname}
        />
      </InputGroup>
      <InputGroup className="mb-3 not-required__field-reg">
        <IMaskInput
          mask="+{7} (000) 000-00-00"
          placeholder={t("auth.phone")}
          id="phone"
          className="form-control"
          onChange={(e) => setPhone(e.target.value)}
          value={phone}
        />
      </InputGroup>
      <InputGroup className="mb-3">
        <span className="required__field required__field-reg">*</span>
        <Form.Control
          placeholder={t("auth.password")}
          type="password"
          minLength={5}
          onChange={(e) => setPassword(e.target.value)}
          value={password}
        />
      </InputGroup>
      <div className="auth-page__action">
        <Button onClick={submit} variant="primary">
          Войти
        </Button>
      </div>
    </div>
  );
};

export default Registration;
