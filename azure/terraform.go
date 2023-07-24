package azure

import (
	"fmt"
	"os/exec"

	"github.com/vision-cli/common/execute"
)

func CallTerraformInit(executor execute.Executor) error {
	fmt.Println("executing make init (terraform)")
	c := exec.Command("make", "init")

	if err := executor.Errors(c, "./azure/_templates/az/tf/", "inititalise Terraform"); err != nil {
		return err
	}

	fmt.Println("make init (terraform) succeeded")

	return nil
}

func CallTerrformPlan(executor execute.Executor) error {
	fmt.Println("executing make plan (terraform)")
	c := exec.Command("make", "plan")

	if err := executor.Errors(c, "./azure/_templates/az/tf/", "plan Terraform"); err != nil {
		return err
	}
	
	println("make plan (terraform) succeeded")
	return nil
}

func CallTerraformApply(executor execute.Executor) error {

	fmt.Println("executing make apply (terraform)")
	c := exec.Command("make", "apply")

	if err := executor.Errors(c, "./azure/_templates/az/tf/", "apply Terraform"); err != nil {
		return err
	}
	
	println("make plan (terraform) succeeded")
	return nil
}