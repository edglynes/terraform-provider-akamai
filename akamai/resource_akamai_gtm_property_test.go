package akamai

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiGTMPropertyConfig = fmt.Sprintf(`
provider "akamai" {
  gtm_section = "gtm"
}

locals {
  	domain = "%s"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_gtm_domain" "test_domain" {
        name = "${local.domain}"
        type = "weighted"
        //contract = "${local.contract}"
	contract = "${data.akamai_contract.contract.id}"
	comment =  "This is a test zone"
	//group     = "${local.group}"
	group = "${data.akamai_group.group.id}"
	wait_on_complete = true
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    nickname = "test_datacenter"
    wait_on_complete = true
    default_load_object = [{
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }]
    depends_on = [
         "akamai_gtm_domain.test_domain"
    ]
}   

resource "akamai_gtm_property" "test_property" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    name = "test_property"
    type = "weighted-round-robin"
    score_aggregation_type = "median"
    handout_limit = 5
    handout_mode = "normal"
    traffic_targets = [{
        datacenter_id = "${akamai_gtm_datacenter.test_datacenter.datacenter_id}"
        enabled = true
        weight = 100
        servers = ["1.2.3.4"]
        // optional
        name = ""
        handout_cname = ""
        }]
    liveness_tests = [{
        name = "lt1"
        test_interval = 30
        test_object_protocol = "HTTP"
        test_timeout = 20
        answer_required = false
        disable_nonstandard_port_warning = false
        error_penalty = 0
        host_header = ""
        http_error3xx = false
        http_error4xx = false
        http_error5xx = false
        peer_certificate_verification = false
        recursion_requested = false
        request_string = ""
        resource_type = ""
        response_string = ""
        ssl_client_certificate = ""
        ssl_client_private_key = ""
        test_object = "junk"
        test_object_password = ""
        test_object_port = 1
        test_object_username = ""
        timeout_penalty = 0
	}]
    depends_on = [
        "akamai_gtm_domain.test_domain",
	"akamai_gtm_datacenter.test_datacenter"
    ]
}
`, gtm_test_domain)

var testAccAkamaiGTMPropertyUpdateConfig = fmt.Sprintf(`
provider "akamai" {
  gtm_section = "gtm"
} 

locals {
        domain = "%s"
}       

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_gtm_domain" "test_domain" {
        name = "${local.domain}"
        type = "weighted"
        contract = "${data.akamai_contract.contract.id}"
        comment =  "This is a test zone"
        group   = "${data.akamai_group.group.id}"
        wait_on_complete = true
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    nickname = "test_datacenter"
    wait_on_complete = true
    default_load_object = [{
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }]
    depends_on = [
         "akamai_gtm_domain.test_domain"
    ]
}

resource "akamai_gtm_property" "test_property" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    name = "test_property"
    type = "weighted-round-robin"
    score_aggregation_type = "median"
    handout_limit = 6
    handout_mode = "normal"
    traffic_targets = [{
        datacenter_id = "${akamai_gtm_datacenter.test_datacenter.datacenter_id}"
        enabled = true
        weight = 100
        servers = ["1.2.3.4"]
        // optional
        name = ""
        handout_cname = ""
        }]
    liveness_tests = [{
        name = "lt1"
        test_interval = 30
        test_object_protocol = "HTTP"
        test_timeout = 20
        answer_required = false
        disable_nonstandard_port_warning = false
        error_penalty = 0
        host_header = ""
        http_error3xx = false
        http_error4xx = false
        http_error5xx = false
        peer_certificate_verification = false
        recursion_requested = false
        request_string = ""
        resource_type = ""
        response_string = ""
        ssl_client_certificate = ""
        ssl_client_private_key = ""
        test_object = "junk"
        test_object_password = ""
        test_object_port = 1
        test_object_username = ""
        timeout_penalty = 0
        }]
    depends_on = [
        "akamai_gtm_domain.test_domain",
	"akamai_gtm_datacenter.test_datacenter"
    ]    
}   
`, gtm_test_domain)

func TestAccAkamaiGTMProperty_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMPropertyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMPropertyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMPropertyExists,
					resource.TestCheckResourceAttr("akamai_gtm_property.test_property", "handout_limit", "5"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMProperty_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMPropertyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMPropertyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMPropertyExists,
					resource.TestCheckResourceAttr("akamai_gtm_property.test_property", "handout_limit", "5"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMPropertyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMPropertyExists,
					resource.TestCheckResourceAttr("akamai_gtm_property.test_property", "handout_limit", "6"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAkamaiGTMPropertyDestroy(s *terraform.State) error {

	// The API doesn't currently support Property Delete
	return nil

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_property" {
			continue
		}

		prop, dom, err := parseStringID(rs.Primary.ID)
		p, err := gtm.GetProperty(prop, dom)
		if p == nil {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMV1_3] Deleting test property [%v]", prop)
		_, err = p.Delete(prop)
		if err != nil {
			return fmt.Errorf("property was not deleted %s. Error: %s", rs.Primary.ID, err.Error())
		}
	}
	return nil
}

func parseStringID(id string) (string, string, error) {
	idComp := strings.Split(id, ":")
	if len(idComp) < 2 {
		return "", "", errors.New("Invalid Property ID")
	}

	return idComp[1], idComp[0], nil

}

func testAccCheckAkamaiGTMPropertyExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_property" {
			continue
		}

		prop, dom, err := parseStringID(rs.Primary.ID)
		_, err = gtm.GetProperty(prop, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
