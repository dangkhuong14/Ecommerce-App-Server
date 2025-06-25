package consumer

import (
	"context"
	"ecommerce/common"
	"ecommerce/common/asyncjob"
	"ecommerce/common/pubsub"
	"ecommerce/component"
	"ecommerce/module/image"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc"
)

var SetImgActiveAfterAvtChangeCmd = &cobra.Command{
	Use:   "SetImgActiveAfterAvtChange",
	Short: "Start consumer: SetImgActiveAfterAvtChange",
	Run: func(cmd *cobra.Command, args []string) {
		// create new service
		newService := sctx.NewServiceContext(sctx.WithName("SetImgActiveAfterAvtChange"),
			sctx.WithComponent(gormc.NewGormDB(common.KeyGormComponent, "")),
			sctx.WithComponent(component.NewNATSComponent(common.KeyNatsComponent)),
		)

		if err := newService.Load(); err != nil {
			log.Fatalln("Error loading service: ", err)
		}

		// get nats service (pubsub)
		ps := newService.MustGet(common.KeyNatsComponent).(pubsub.PubSub)

		// subscribe to user change avatar topic
		ctx := context.Background()
		ch, _ := ps.Subscribe(ctx, common.TopicUserChangedAvatar)

		// get image's repo
		gorm := newService.MustGet(common.KeyGormComponent).(common.GormCompContext)
		imgRepo := image.NewImageRepo(gorm.GetDB())

		// loop message channel
		for msg := range ch {
			go func() {
				defer common.Recover()
				mapData := msg.Data()
				imgID := common.UUID(uuid.MustParse(mapData["img_id"].(string)))
				// update image's status using async job
				// create job
				job := asyncjob.NewJob(func(ctx context.Context) error {
					return imgRepo.SetImageStatusActivated(ctx, imgID)
				}, asyncjob.WithName("SetImageStatusActivated"))

				// create group
				group := asyncjob.NewGroup(true, job)
				group.Run(ctx)
			}()
		}
	},
}
