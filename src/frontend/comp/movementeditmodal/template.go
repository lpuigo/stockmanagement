package movementeditmodal

const template string = `<el-dialog
        :before-close="HideWithControl"
        :visible.sync="Visible" :close-on-click-modal="false"
        width="70vw" top="10vh"
>
    <!-- 
        Modal Title
    -->
    <div slot="title">
		<h2 style="margin: 0 10px" v-if="current_movement">
			<i class="far fa-edit icon--left"></i>Edition de <span class="batec-header-text">{{FormatMovementType(edited_movement.Type)}} du {{FormatDate(edited_movement.Date)}}</span>
		</h2>
    </div>
    <!-- 
        Modal Body
    -->
	<el-tabs tab-position="top" v-model="EditMode" style="margin: 0px 15px;">
		<!-- ====================	Movement Edition  ============================== -->
		<el-tab-pane :label="FormatMovementType(current_movement.Type)" name="acc" lazy=true style="padding: 0px 0px">
			<el-container style="margin: 0px;height: 60vh">
				<el-header>
					<!--	Actor & Responsible -->
					<el-row :gutter="5" type="flex" align="middle" class="spaced">
						<el-col :span="2" class="align-right">Date:</el-col>
						<el-col :span="10">
							<el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
								style="width: 100%" type="date"
								v-model="current_movement.Date"
								value-format="yyyy-MM-dd"
								:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
								@change="UpdateDate"
							></el-date-picker>
						</el-col>
						<el-col :span="2" class="align-right">Acteur:</el-col>
						<el-col :span="10">
							<el-input v-model="current_movement.Actor" size="mini""></el-input>
						</el-col>
					</el-row>
			
					<!--	Date  -->
					<el-row v-if="IsWorksiteRelated()" :gutter="5" type="flex" align="middle" class="spaced">
						<el-col :span="2" class="align-right">Chantier:</el-col>
						<el-col :span="10">
							<el-select v-model="current_movement.WorksiteId"  
								filterable size="mini" style="width: 100%"
								@change="UpdateWorksite">
								<el-option
									v-for="item in GetActiveWorksites()" :key="item.value"
									:value="item.value"
									:label="item.label"
								></el-option>
							</el-select>
						</el-col>
						<el-col :span="2" class="align-right">Responsable:</el-col>
						<el-col :span="10">
							<el-input v-model="current_movement.Responsible" size="mini""></el-input>
						</el-col>
					</el-row>
				</el-header>
				<el-main>
					<articleflow-table
                        v-model="current_movement.ArticleFlows"
                        :articles="StockArticles"
                        :user="User"
					></articleflow-table>
				</el-main>				
			</el-container>
		</el-tab-pane>

		<!-- ====================	Rental Stays Edition  ============================== -->
		<el-tab-pane label="Mouvement" name="movement" lazy=true style="padding: 0px 20px;height: 40vh;overflow-y: auto">
			<pre>{{current_movement}}</pre>
		</el-tab-pane>

		<!-- ====================	Rental Stays Edition  ============================== -->
		<el-tab-pane label="Stock" name="stock" lazy=true style="padding: 0px 20px;height: 40vh;overflow-y: auto">
			<pre>{{stock}}</pre>
		</el-tab-pane>

		<!-- ====================	Rental Stays Edition  ============================== -->
		<el-tab-pane label="Articles" name="articles" lazy=true style="padding: 0px 20px;height: 40vh;overflow-y: auto">
			<pre>{{articles.Articles}}</pre>
		</el-tab-pane>

	</el-tabs>

    <!-- 
        Modal Footer Action Bar
    -->
    <div slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-tooltip :open-delay="500" effect="light">
					<div slot="content">Annuler les changements</div>
					<el-button :disabled="!hasChanged" @click="UndoChange" icon="fas fa-undo-alt" plain size="mini"
                               type="info"></el-button>
				</el-tooltip>
				
				<el-button @click="CancelChange" size="mini">Fermer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Valider</el-button>
			</el-col>
		</el-row>
	</div>
</el-dialog>`
