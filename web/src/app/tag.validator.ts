import {
  AbstractControl,
  FormArray,
  NG_VALIDATORS,
  ValidationErrors,
  Validator,
} from '@angular/forms';
import { Directive, forwardRef } from '@angular/core';
import { debounceTime, ReplaySubject } from 'rxjs';
import { ValidateRequest, ValidatorService } from './validator.service';

@Directive({
  providers: [
    {
      multi: true,
      provide: NG_VALIDATORS,
      useExisting: forwardRef(() => TagFormValidator),
    },
  ],
  selector: '[appTagFormValidator]',
})
export class TagFormValidator implements Validator {
  private debounceSubject = new ReplaySubject<() => void>(1);

  constructor(private validator: ValidatorService) {
    this.debounceSubject.pipe(debounceTime(500)).subscribe((func) => func());
  }

  validate(control: AbstractControl): ValidationErrors | null {
    this.debounceSubject.next(() => {
      const tags = control.get('tags') as FormArray;
      let toValidate: ValidateRequest = {
        data: {},
      };

      tags.controls.forEach((tag) => {
        tag.get('value')?.setErrors(null);

        const key = tag.value['key'];

        if (key !== undefined && key !== null && key !== '') {
          toValidate.data[key] = tag.value['value'];
        }
      });

      this.validator.validateTag(toValidate).subscribe((response) => {
        Object.entries(response.data).forEach(([key, value]) => {
          let control = tags.controls.find((tag) => tag.value['key'] === key);
          if (control !== undefined) {
            control.get('value')?.setErrors({
              serverMessage: value,
            });
          }
        });

        control.markAllAsTouched();
      });
    });
    return null;
  }
}
