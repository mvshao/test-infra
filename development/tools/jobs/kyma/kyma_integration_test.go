package kyma_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/test-infra/development/tools/jobs/tester"
	"github.com/kyma-project/test-infra/development/tools/jobs/tester/preset"
)

func TestKymaIntegrationJobsPresubmit(t *testing.T) {
	tests := map[string]struct {
		givenJobName            string
		expPresets              []preset.Preset
		expRunIfChangedRegex    string
		expRunIfChangedPaths    []string
		expNotRunIfChangedPaths []string
	}{
		"Should contain the kyma-integration k3d job": {
			givenJobName: "pre-main-kyma-integration-k3d",

			expPresets: []preset.Preset{
				preset.GCProjectEnv, preset.KymaGuardBotGithubToken, preset.BuildPr, "preset-sa-vm-kyma-integration",
			},

			expRunIfChangedRegex: "^((tests/fast-integration\\S+|resources\\S+|installation\\S+|tools/kyma-installer\\S+)(\\.[^.][^.][^.]+$|\\.[^.][^dD]$|\\.[^mM][^.]$|\\.[^.]$|/[^.]+$))",
			expRunIfChangedPaths: []string{
				"resources/values.yaml",
				"installation/file.yaml",
			},
			expNotRunIfChangedPaths: []string{
				"installation/README.md",
				"installation/test/test/README.MD",
			},
		},
		"Should contain the kyma-integration k3d with central Application Connectivity job": {
			givenJobName: "pre-main-kyma-integration-k3d-central-app-connectivity",

			expPresets: []preset.Preset{
				preset.GCProjectEnv, preset.KymaGuardBotGithubToken, preset.BuildPr, "preset-sa-vm-kyma-integration", "preset-kyma-integration-central-app-connectivity-enabled",
			},

			expRunIfChangedRegex: "^((tests/fast-integration\\S+|resources\\S+|installation\\S+|tools/kyma-installer\\S+)(\\.[^.][^.][^.]+$|\\.[^.][^dD]$|\\.[^mM][^.]$|\\.[^.]$|/[^.]+$))",
			expRunIfChangedPaths: []string{
				"resources/values.yaml",
				"installation/file.yaml",
			},
			expNotRunIfChangedPaths: []string{
				"installation/README.md",
				"installation/test/test/README.MD",
			},
		},
		"Should contain the kyma-integration k3d with central Application Connectivity and Compass job": {
			givenJobName: "pre-main-kyma-integration-k3d-central-app-connectivity-compass",

			expPresets: []preset.Preset{
				preset.GCProjectEnv, preset.KymaGuardBotGithubToken, preset.BuildPr, "preset-sa-vm-kyma-integration", "preset-kyma-integration-central-app-connectivity-enabled", "preset-kyma-integration-compass-dev", "preset-kyma-integration-compass-enabled",
			},

			expRunIfChangedRegex: "^((tests/fast-integration\\S+|resources\\S+|installation\\S+|tools/kyma-installer\\S+)(\\.[^.][^.][^.]+$|\\.[^.][^dD]$|\\.[^mM][^.]$|\\.[^.]$|/[^.]+$))",
			expRunIfChangedPaths: []string{
				"resources/values.yaml",
				"installation/file.yaml",
			},
			expNotRunIfChangedPaths: []string{
				"installation/README.md",
				"installation/test/test/README.MD",
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			jobConfig, err := tester.ReadJobConfig("./../../../../prow/jobs/kyma/kyma-integration.yaml")
			require.NoError(t, err)

			// when
			actualJob := tester.FindPresubmitJobByNameAndBranch(jobConfig.AllStaticPresubmits([]string{"kyma-project/kyma"}), tc.givenJobName, "main")
			require.NotNil(t, actualJob)

			// then
			// the common expectation
			assert.Equal(t, "github.com/kyma-project/kyma", actualJob.PathAlias)
			assert.Equal(t, tc.expRunIfChangedRegex, actualJob.RunIfChanged)

			assert.False(t, actualJob.SkipReport)
			assert.Equal(t, 10, actualJob.MaxConcurrency)
			tester.AssertThatHasExtraRefTestInfra(t, actualJob.JobBase.UtilityConfig, "main")
			tester.AssertThatSpecifiesResourceRequests(t, actualJob.JobBase)

			// the job specific expectation
			tester.AssertThatHasPresets(t, actualJob.JobBase, tc.expPresets...)
			for _, path := range tc.expRunIfChangedPaths {
				assert.True(t, tester.IfPresubmitShouldRunAgainstChanges(*actualJob, true, path))
			}
			for _, path := range tc.expNotRunIfChangedPaths {
				assert.False(t, tester.IfPresubmitShouldRunAgainstChanges(*actualJob, true, path))
			}
		})
	}
}

func TestKymaIntegrationJobsPostsubmit(t *testing.T) {
	tests := map[string]struct {
		givenJobName string
		expPresets   []preset.Preset
	}{

		"Should contain the kyma-integration-k3d job": {
			givenJobName: "post-main-kyma-integration-k3d",

			expPresets: []preset.Preset{
				preset.GCProjectEnv, preset.KymaGuardBotGithubToken, "preset-sa-vm-kyma-integration",
			},
		},
		"Should contain the kyma-integration k3d with central Application Connectivity job": {
			givenJobName: "post-main-kyma-integration-k3d-central-app-connectivity",

			expPresets: []preset.Preset{
				preset.GCProjectEnv, preset.KymaGuardBotGithubToken, "preset-sa-vm-kyma-integration", "preset-kyma-integration-central-app-connectivity-enabled",
			},
		},
		"Should contain the kyma-integration k3d with central Application Connectivity and Compass job": {
			givenJobName: "post-main-kyma-integration-k3d-central-app-connectivity-compass",

			expPresets: []preset.Preset{
				preset.GCProjectEnv, preset.KymaGuardBotGithubToken, "preset-sa-vm-kyma-integration", "preset-kyma-integration-central-app-connectivity-enabled", "preset-kyma-integration-compass-dev", "preset-kyma-integration-compass-enabled",
			},
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			jobConfig, err := tester.ReadJobConfig("./../../../../prow/jobs/kyma/kyma-integration.yaml")
			require.NoError(t, err)

			// when
			actualJob := tester.FindPostsubmitJobByNameAndBranch(jobConfig.AllStaticPostsubmits([]string{"kyma-project/kyma"}), tc.givenJobName, "main")
			require.NotNil(t, actualJob)

			// then
			// the common expectation
			assert.Equal(t, []string{"^master$", "^main$"}, actualJob.Branches)
			assert.Equal(t, 10, actualJob.MaxConcurrency)
			assert.Equal(t, "", actualJob.RunIfChanged)

			assert.Equal(t, "github.com/kyma-project/kyma", actualJob.PathAlias)
			tester.AssertThatHasExtraRefTestInfra(t, actualJob.JobBase.UtilityConfig, "main")
			tester.AssertThatSpecifiesResourceRequests(t, actualJob.JobBase)

			// the job specific expectation
			tester.AssertThatHasPresets(t, actualJob.JobBase, tc.expPresets...)
		})
	}
}

func TestKymaIntegrationJobPeriodics(t *testing.T) {
	// WHEN
	jobConfig, err := tester.ReadJobConfig("./../../../../prow/jobs/kyma/kyma-integration.yaml")
	// THEN
	require.NoError(t, err)

	periodics := jobConfig.AllPeriodics()
	assert.Len(t, periodics, 18)

	expName := "kyma-upgrade-k3d-kyma-to-kyma2"
	kymaUpgradePeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, kymaUpgradePeriodic)
	assert.Equal(t, expName, kymaUpgradePeriodic.Name)

	assert.Equal(t, "0 0 6-18/2 ? * 1-5", kymaUpgradePeriodic.Cron)
	tester.AssertThatHasPresets(t, kymaUpgradePeriodic.JobBase, preset.GCProjectEnv, preset.KymaGuardBotGithubToken, "preset-sa-vm-kyma-integration")
	assert.Equal(t, tester.ImageKymaIntegrationLatest, kymaUpgradePeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/provision-vm-and-start-kyma-upgrade-k3d.sh"}, kymaUpgradePeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string(nil), kymaUpgradePeriodic.Spec.Containers[0].Args)
	tester.AssertThatContainerHasEnv(t, kymaUpgradePeriodic.Spec.Containers[0], "KYMA_PROJECT_DIR", ".")
	tester.AssertThatContainerHasEnv(t, kymaUpgradePeriodic.Spec.Containers[0], "KYMA_MAJOR_VERSION", "1")
	tester.AssertThatSpecifiesResourceRequests(t, kymaUpgradePeriodic.JobBase)

	expName = "kyma-upgrade-k3d-kyma2-to-main"
	kymaUpgradePeriodic = tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, kymaUpgradePeriodic)
	assert.Equal(t, expName, kymaUpgradePeriodic.Name)

	assert.Equal(t, "0 0 6-18/2 ? * 1-5", kymaUpgradePeriodic.Cron)
	tester.AssertThatHasPresets(t, kymaUpgradePeriodic.JobBase, preset.GCProjectEnv, preset.KymaGuardBotGithubToken, "preset-sa-vm-kyma-integration")
	assert.Equal(t, tester.ImageKymaIntegrationLatest, kymaUpgradePeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/provision-vm-and-start-kyma-upgrade-k3d.sh"}, kymaUpgradePeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string(nil), kymaUpgradePeriodic.Spec.Containers[0].Args)
	tester.AssertThatContainerHasEnv(t, kymaUpgradePeriodic.Spec.Containers[0], "KYMA_PROJECT_DIR", ".")
	tester.AssertThatContainerHasEnv(t, kymaUpgradePeriodic.Spec.Containers[0], "KYMA_MAJOR_VERSION", "2")
	tester.AssertThatSpecifiesResourceRequests(t, kymaUpgradePeriodic.JobBase)

	expName = "orphaned-disks-cleaner"
	disksCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, disksCleanerPeriodic)
	assert.Equal(t, expName, disksCleanerPeriodic.Name)

	assert.Equal(t, "30 * * * *", disksCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, disksCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, disksCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, disksCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, disksCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/diskscollector -project=${CLOUDSDK_CORE_PROJECT} -dryRun=false -diskNameRegex='^gke-'"}, disksCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, disksCleanerPeriodic.JobBase)

	expName = "orphaned-ips-cleaner"
	addressesCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, addressesCleanerPeriodic)
	assert.Equal(t, expName, addressesCleanerPeriodic.Name)

	assert.Equal(t, "0 1 * * *", addressesCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, addressesCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, addressesCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, addressesCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, addressesCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/ipcleaner -project=${CLOUDSDK_CORE_PROJECT} -dry-run=false -ip-exclude-name-regex='^nightly|nightly-124|weekly|weekly-124|nat-auto-ip|nightly-20'"}, addressesCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, addressesCleanerPeriodic.JobBase)

	expName = "orphaned-az-storage-accounts-cleaner"
	orphanedAZStorageAccountsCleaner := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, orphanedAZStorageAccountsCleaner)
	assert.Equal(t, expName, orphanedAZStorageAccountsCleaner.Name)

	assert.Equal(t, "00 00 * * *", orphanedAZStorageAccountsCleaner.Cron)
	tester.AssertThatHasPresets(t, orphanedAZStorageAccountsCleaner.JobBase, "preset-az-kyma-prow-credentials")
	tester.AssertThatHasExtraRepoRefCustom(t, orphanedAZStorageAccountsCleaner.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageKymaIntegrationLatest, orphanedAZStorageAccountsCleaner.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/cluster-integration/minio/azure-cleaner.sh"}, orphanedAZStorageAccountsCleaner.Spec.Containers[0].Command)
	tester.AssertThatSpecifiesResourceRequests(t, orphanedAZStorageAccountsCleaner.JobBase)

	expName = "orphaned-assetstore-gcp-bucket-cleaner"
	assetstoreGcpBucketCleaner := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, assetstoreGcpBucketCleaner)
	assert.Equal(t, expName, assetstoreGcpBucketCleaner.Name)

	assert.Equal(t, "00 00 * * *", assetstoreGcpBucketCleaner.Cron)
	tester.AssertThatHasPresets(t, assetstoreGcpBucketCleaner.JobBase, preset.GCProjectEnv, preset.SaProwJobResourceCleaner)
	tester.AssertThatHasExtraRepoRefCustom(t, assetstoreGcpBucketCleaner.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, assetstoreGcpBucketCleaner.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, assetstoreGcpBucketCleaner.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "prow/scripts/assetstore-gcp-bucket-cleaner.sh -project=${CLOUDSDK_CORE_PROJECT}"}, assetstoreGcpBucketCleaner.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, assetstoreGcpBucketCleaner.JobBase)

	expName = "orphaned-clusters-cleaner"
	clustersCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, clustersCleanerPeriodic)
	assert.Equal(t, expName, clustersCleanerPeriodic.Name)

	assert.Equal(t, "0 * * * *", clustersCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, clustersCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, clustersCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, clustersCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, clustersCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/clusterscollector -project=${CLOUDSDK_CORE_PROJECT} -dryRun=false -excluded-clusters=kyma-prow,workload-kyma-prow,nightly,weekly,nightly-124,weekly-124,nightly-20"}, clustersCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, clustersCleanerPeriodic.JobBase)

	expName = "orphaned-vms-cleaner"
	vmsCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, vmsCleanerPeriodic)
	assert.Equal(t, expName, vmsCleanerPeriodic.Name)

	assert.Equal(t, "15,45 * * * *", vmsCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, vmsCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, vmsCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, vmsCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, vmsCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/vmscollector -project=${CLOUDSDK_CORE_PROJECT} -vmNameRegexp='.*-integration-test-.*|busola-integration-test-.*' -jobLabelRegexp='.*-integration$|busola-integration-test-k3s' -dryRun=false"}, vmsCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, vmsCleanerPeriodic.JobBase)

	expName = "orphaned-loadbalancer-cleaner"
	loadbalancerCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, loadbalancerCleanerPeriodic)
	assert.Equal(t, expName, loadbalancerCleanerPeriodic.Name)

	assert.Equal(t, "15 * * * *", loadbalancerCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, loadbalancerCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, loadbalancerCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, loadbalancerCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, loadbalancerCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/orphanremover -project=${CLOUDSDK_CORE_PROJECT} -dryRun=false"}, loadbalancerCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, loadbalancerCleanerPeriodic.JobBase)

	// expName = "firewall-cleaner"
	// firewallCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	// require.NotNil(t, firewallCleanerPeriodic)
	// assert.Equal(t, expName, firewallCleanerPeriodic.Name)
	//
	// assert.Equal(t, "45 */4 * * 1-5", firewallCleanerPeriodic.Cron)
	// tester.AssertThatHasPresets(t, firewallCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	// tester.AssertThatHasExtraRepoRefCustom(t, firewallCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	// assert.Equal(t, tester.ImageProwToolsLatest, firewallCleanerPeriodic.Spec.Containers[0].Image)
	// assert.Equal(t, []string{"bash"}, firewallCleanerPeriodic.Spec.Containers[0].Command)
	// assert.Equal(t, []string{"-c", "/prow-tools/firewallcleaner -project=${CLOUDSDK_CORE_PROJECT} -dryRun=false"}, firewallCleanerPeriodic.Spec.Containers[0].Args)
	// tester.AssertThatSpecifiesResourceRequests(t, firewallCleanerPeriodic.JobBase)

	expName = "orphaned-dns-cleaner"
	dnsCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, dnsCleanerPeriodic)
	assert.Equal(t, expName, dnsCleanerPeriodic.Name)

	assert.Equal(t, "30 * * * *", dnsCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, dnsCleanerPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, dnsCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, dnsCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, dnsCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/dnscollector -project=${CLOUDSDK_CORE_PROJECT} -dnsZone=${CLOUDSDK_DNS_ZONE_NAME} -ageInHours=2 -regions=${CLOUDSDK_COMPUTE_REGION} -dryRun=false"}, dnsCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, dnsCleanerPeriodic.JobBase)

	expName = "gcr-cleaner-prow-workloads"
	gcrCleanerPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, gcrCleanerPeriodic)
	assert.Equal(t, expName, gcrCleanerPeriodic.Name)

	assert.Equal(t, "0 1 * * *", gcrCleanerPeriodic.Cron)
	tester.AssertThatHasPresets(t, gcrCleanerPeriodic.JobBase, preset.SaGKEKymaIntegration)
	tester.AssertThatHasExtraRepoRefCustom(t, gcrCleanerPeriodic.JobBase.UtilityConfig, []string{"test-infra"}, []string{"main"})
	assert.Equal(t, tester.ImageProwToolsLatest, gcrCleanerPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"bash"}, gcrCleanerPeriodic.Spec.Containers[0].Command)
	assert.Equal(t, []string{"-c", "/prow-tools/gcrcleaner --repository=eu.gcr.io/sap-kyma-prow-workloads --age-in-hours=168 --dry-run=false"}, gcrCleanerPeriodic.Spec.Containers[0].Args)
	tester.AssertThatSpecifiesResourceRequests(t, gcrCleanerPeriodic.JobBase)

	expName = "github-stats"
	githubStatsPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, githubStatsPeriodic)
	assert.Equal(t, expName, githubStatsPeriodic.Name)

	assert.Equal(t, "0 6 * * *", githubStatsPeriodic.Cron)
	assert.Equal(t, tester.ImageProwToolsLatest, githubStatsPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/github-stats.sh"}, githubStatsPeriodic.Spec.Containers[0].Command)
	tester.AssertThatSpecifiesResourceRequests(t, githubStatsPeriodic.JobBase)

	expName = "github-issues"
	githubIssuesPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, githubIssuesPeriodic)
	assert.Equal(t, expName, githubIssuesPeriodic.Name)

	assert.Equal(t, "0 6 * * *", githubIssuesPeriodic.Cron)
	assert.Equal(t, tester.ImageProwToolsLatest, githubIssuesPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/github-issues.sh"}, githubIssuesPeriodic.Spec.Containers[0].Command)
	tester.AssertThatSpecifiesResourceRequests(t, githubIssuesPeriodic.JobBase)

	expName = "kyma-gke-nightly-fast-integration"
	nightlyFastIntegrationPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	require.NotNil(t, nightlyFastIntegrationPeriodic)
	assert.Equal(t, expName, nightlyFastIntegrationPeriodic.Name)

	assert.Equal(t, "0 0-1,4-23 * * 1-5", nightlyFastIntegrationPeriodic.Cron)
	tester.AssertThatHasPresets(t, nightlyFastIntegrationPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration, "preset-gc-compute-envs")
	tester.AssertThatHasExtraRepoRefCustom(t, nightlyFastIntegrationPeriodic.JobBase.UtilityConfig, []string{"test-infra", "kyma"}, []string{"main", "main"})
	assert.Equal(t, tester.ImageKymaIntegrationLatest, nightlyFastIntegrationPeriodic.Spec.Containers[0].Image)
	assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/cluster-integration/fast-integration-test.sh"}, nightlyFastIntegrationPeriodic.Spec.Containers[0].Command)
	tester.AssertThatSpecifiesResourceRequests(t, nightlyFastIntegrationPeriodic.JobBase)
	assert.Len(t, nightlyFastIntegrationPeriodic.Spec.Containers[0].Env, 5)
	tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "PROVISION_REGIONAL_CLUSTER", "true")
	tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "INPUT_CLUSTER_NAME", "nightly")
	tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "CLUSTER_PROVIDER", "gcp")
	tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "CLOUDSDK_COMPUTE_ZONE", "europe-west4-b")
	tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "KYMA_MAJOR_VERSION", "2")

	// these jobs are disabled due to failing tests, until Kyma 2.0 long lasting clusters are introduced
	// expName = "kyma-gke-weekly-fast-integration"
	// weeklyFastIntegrationPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	// require.NotNil(t, weeklyFastIntegrationPeriodic)
	// assert.Equal(t, expName, weeklyFastIntegrationPeriodic.Name)

	// assert.Equal(t, "0 0-4,6-23 * * 1-5", weeklyFastIntegrationPeriodic.Cron)
	// tester.AssertThatHasPresets(t, weeklyFastIntegrationPeriodic.JobBase, preset.GCProjectEnv, preset.SaGKEKymaIntegration, "preset-gc-compute-envs")
	// tester.AssertThatHasExtraRepoRefCustom(t, weeklyFastIntegrationPeriodic.JobBase.UtilityConfig, []string{"test-infra", "kyma"}, []string{"main", "main"})
	// assert.Equal(t, tester.ImageKymaIntegrationLatest, weeklyFastIntegrationPeriodic.Spec.Containers[0].Image)
	// assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/cluster-integration/fast-integration-test.sh"}, weeklyFastIntegrationPeriodic.Spec.Containers[0].Command)
	// tester.AssertThatSpecifiesResourceRequests(t, weeklyFastIntegrationPeriodic.JobBase)
	// assert.Len(t, weeklyFastIntegrationPeriodic.Spec.Containers[0].Env, 5)
	// tester.AssertThatContainerHasEnv(t, weeklyFastIntegrationPeriodic.Spec.Containers[0], "PROVISION_REGIONAL_CLUSTER", "true")
	// tester.AssertThatContainerHasEnv(t, weeklyFastIntegrationPeriodic.Spec.Containers[0], "INPUT_CLUSTER_NAME", "weekly-124")
	// tester.AssertThatContainerHasEnv(t, weeklyFastIntegrationPeriodic.Spec.Containers[0], "CLUSTER_PROVIDER", "gcp")
	// tester.AssertThatContainerHasEnv(t, weeklyFastIntegrationPeriodic.Spec.Containers[0], "CLOUDSDK_COMPUTE_ZONE", "europe-west4-b")
	// tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "KYMA_MAJOR_VERSION", "1")

	// expName = "kyma-aks-nightly-fast-integration"
	// nightlyAksFastIntegrationPeriodic := tester.FindPeriodicJobByName(periodics, expName)
	// require.NotNil(t, nightlyAksFastIntegrationPeriodic)
	// assert.Equal(t, expName, nightlyAksFastIntegrationPeriodic.Name)

	// assert.Equal(t, "0 0-2,4-23 * * 1-5", nightlyAksFastIntegrationPeriodic.Cron)
	// tester.AssertThatHasPresets(t, nightlyAksFastIntegrationPeriodic.JobBase, preset.SaGKEKymaIntegration, "preset-az-kyma-prow-credentials")
	// tester.AssertThatHasExtraRepoRefCustom(t, nightlyAksFastIntegrationPeriodic.JobBase.UtilityConfig, []string{"test-infra", "kyma"}, []string{"main", "main"})
	// assert.Equal(t, tester.ImageKymaIntegrationLatest, nightlyAksFastIntegrationPeriodic.Spec.Containers[0].Image)
	// assert.Equal(t, []string{"/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/cluster-integration/fast-integration-test.sh"}, nightlyAksFastIntegrationPeriodic.Spec.Containers[0].Command)
	// tester.AssertThatSpecifiesResourceRequests(t, nightlyAksFastIntegrationPeriodic.JobBase)
	// assert.Len(t, nightlyAksFastIntegrationPeriodic.Spec.Containers[0].Env, 4)
	// tester.AssertThatContainerHasEnv(t, nightlyAksFastIntegrationPeriodic.Spec.Containers[0], "RS_GROUP", "kyma-nightly-aks")
	// tester.AssertThatContainerHasEnv(t, nightlyAksFastIntegrationPeriodic.Spec.Containers[0], "INPUT_CLUSTER_NAME", "nightly-aks-124")
	// tester.AssertThatContainerHasEnv(t, nightlyAksFastIntegrationPeriodic.Spec.Containers[0], "CLUSTER_PROVIDER", "azure")
	// tester.AssertThatContainerHasEnv(t, nightlyFastIntegrationPeriodic.Spec.Containers[0], "KYMA_MAJOR_VERSION", "1")
}
