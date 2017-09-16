import { Component, OnInit, OnDestroy, Inject } from "@angular/core" ;
import { FormControl, Validators } from '@angular/forms';
import { MdDialog, MdDialogRef, MD_DIALOG_DATA } from "@angular/material";
import { UserService, PasswordChange } from "../../services/user";

@Component({
    selector: "settings-general",
    templateUrl: "./general.html",
    styleUrls: ["../common.css"]
})
export class GeneralSettingsComponent implements OnInit {
    firstName: string
    lastName: string
    email: string

    emailFormControl = new FormControl('', [Validators.email]);

    constructor(
        private userService: UserService,
        private dialog: MdDialog,
    ) {}

    ngOnInit(): void {
        this.userService.getCurrentUser().subscribe(
            user => {
                this.firstName = user.firstName;
                this.lastName = user.lastName;
                this.email = user.email;
            },
            error => console.log(error),
        )
    }

    firstNameChange() {
        this.userService.setUserSetting(
            "first-name", this.firstName
        ).subscribe(
            success => {},
            error => console.log(error),
        )
    }

    lastNameChange() {
        this.userService.setUserSetting(
            "last-name", this.lastName
        ).subscribe(
            success => {},
            error => console.log(error),
        )
    }

    emailChange() {
        this.emailFormControl.updateValueAndValidity()

        if (!this.emailFormControl.valid) {
            return
        }

        this.userService.setUserSetting(
            "email", this.email
        ).subscribe(
            success => {
                if (!success) {
                    this.emailFormControl.setErrors({"email": true})
                }
            },
            error => console.log(error),
        )
    }

    changePassword() {
        this.dialog.open(PasswordDialog, {
            width: "250px",
        })
    }
}

@Component({
    templateUrl: "./password-form.html",
    styleUrls: ["../common.css"]
})
export class PasswordDialog {
    current: string
    password: string
    passwordConfirm: string

    currentFormControl = new FormControl('', [Validators.required]);
    passwordFormControl = new FormControl('', [Validators.required]);

    constructor(
        private dialogRef: MdDialogRef<PasswordDialog>,
        private userService: UserService,
        // @Inject(MD_DIALOG_DATA) private data: any,
    ) {}

    save() {
        if (this.password != this.passwordConfirm) {
            this.passwordFormControl.setErrors({"mismatch": true})
            return
        }

        this.currentFormControl.updateValueAndValidity()

        if (!this.currentFormControl.valid) {
            return
        }

        this.userService.changeUserPassword(
            {current: this.current, new: this.password}
        ).subscribe(
            success => {
                if (!success) {
                    this.currentFormControl.setErrors({"auth": true})
                    return
                }

                this.close();
            },
            error => console.log(error),
        )
    }

    close() {
        this.dialogRef.close();
    }
}