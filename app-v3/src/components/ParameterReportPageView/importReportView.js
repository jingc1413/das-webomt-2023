import ReportElement from "./ReportElement.vue"
import ReportForm from "./ReportForm.vue"
import ReportFormElement from "./ReportFormElement.vue"
import ReportFormLayout from "./ReportFormLayout.vue"
import ReportLayout from "./ReportLayout.vue"
import ReportMainPage from "./ReportMainPage.vue"
import ReportParam from "./ReportParam.vue"
import ReportParamGroup from "./ReportParamGroup.vue"
import ReportTableView from "./ReportTableView.vue"
import ReportViewPage from "./ReportViewPage.vue"
import ReportButton from "./ReportButton.vue"
import ReportToolbarForm from "./ReportToolbarForm.vue"

export default function ImportReportViewComponents(appRoot) {

  appRoot.component("report-main-page", ReportMainPage);
  appRoot.component("report-view-page", ReportViewPage);
  appRoot.component("report-layout", ReportLayout);
  appRoot.component("report-element", ReportElement);
  appRoot.component("report-toolbar-form", ReportToolbarForm);
  appRoot.component("report-table", ReportTableView);
  appRoot.component("report-form", ReportForm);
  appRoot.component("report-form-layout", ReportFormLayout);
  appRoot.component("report-form-element", ReportFormElement);
  appRoot.component("report-param-group", ReportParamGroup);
  appRoot.component("report-param", ReportParam);
  appRoot.component("report-button", ReportButton);
}